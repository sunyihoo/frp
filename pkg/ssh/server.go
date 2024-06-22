package ssh

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sunyihoo/frp/pkg/config"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/util/util"
	"net"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"github.com/sunyihoo/frp/pkg/virtual"
)

const (
	RequestTypeForward = "tcpip-forward"
)

type tcpipForward struct {
	Host string
	Port uint32
}

type TunnelServer struct {
	underlyingConn net.Conn
	sshConn        *ssh.ServerConn
	sc             *ssh.ServerConfig
	firstChannel   ssh.Channel

	vc                 *virtual.Client
	peerServerListener *netpkg.InternalListener
	doneCh             chan struct{}
	closeDoneCh        sync.Once
}

func NewTunnelServer(conn net.Conn, sc *ssh.ServerConfig, peerServerListener *netpkg.InternalListener) (*TunnelServer, error) {
	s := &TunnelServer{
		underlyingConn:     conn,
		sc:                 sc,
		peerServerListener: peerServerListener,
		doneCh:             make(chan struct{}),
	}
	return s, nil
}

func (s *TunnelServer) Run() error {
	sshConn, channels, requests, err := ssh.NewServerConn(s.underlyingConn, s.sc)
	if err != nil {
		return nil
	}

	s.sshConn = sshConn

	addr, extraPayload, err := s.waitForwardAddrAndExtraPayload(channels, requests, 3*time.Second)
	if err != nil {
		return err
	}

	clientCfg, pc, helpMessage, err := s.parseClientAndProxyConfigurer(addr, extraPayload)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			s.writeToClient(helpMessage)
			return nil
		}
		s.writeToClient(err.Error())
		return fmt.Errorf("parse flags from ssh client error: %v", err)
	}
	clientCfg.Complete()
	if sshConn.Permissions != nil {
		clientCfg.User = util.EmptyOr(sshConn.Permissions.Extensions["user"], clientCfg.User)
	}
	pc.Complete(clientCfg.User)

	vc, err := virtual.NewClient()

}

func (s *TunnelServer) writeToClient(data string) {
	if s.firstChannel == nil {
		return
	}
	_, _ = s.firstChannel.Write([]byte(data + "\n"))
}
func (s *TunnelServer) waitForwardAddrAndExtraPayload(
	channels <-chan ssh.NewChannel,
	requests <-chan *ssh.Request,
	timeout time.Duration,
) (*tcpipForward, string, error) {
	addrCh := make(chan *tcpipForward, 1)
	extraPayloadCh := make(chan string, 1)

	// get forward address
	go func() {
		addrGot := false
		for req := range requests {
			if req.Type == RequestTypeForward && !addrGot {
				payload := tcpipForward{}
				if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
					return
				}
				addrGot = true
				addrCh <- &payload
			}
			if req.WantReply {
				_ = req.Reply(true, nil)
			}
		}
	}()

	// get extra payload
	go func() {
		for newChannel := range channels {
			// extraPayload will send to extraPayloadCh
			go s.handleNewChannel(newChannel, extraPayloadCh)
		}
	}()

	var (
		addr         *tcpipForward
		extraPayload string
	)

	// todo 学习
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case v := <-addrCh:
			addr = v
		case extra := <-extraPayloadCh:
			extraPayload = extra
		case <-timer.C:
			return nil, "", fmt.Errorf("get addr and extra payload timeout")
		}
		if addr != nil && extraPayload != "" {
			break
		}
	}
	return addr, extraPayload, nil
}

func (s *TunnelServer) parseClientAndProxyConfigurer(_ *tcpipForward, extraPayload string) (*v1.ClientCommonConfig, v1.ProxyConfigurer, string, error) {
	helpMessage := ""
	cmd := &cobra.Command{
		Use:   "ssh v0&{address} [command]",
		Short: "ssh v0&{address} [command]",
		Run:   func(*cobra.Command, []string) {},
	}
	// todo 学习
	cmd.SetGlobalNormalizationFunc(config.WordSepNormalizeFunc)

	args := strings.Split(extraPayload, " ")
	if len(args) < 1 {
		return nil, nil, helpMessage, fmt.Errorf("invalid extra payload")
	}
	proxyType := strings.TrimSpace(args[0])
	supportType := []string{"tcp", "http", "https", "tcpmux", "stcp"}
	if !slices.Contains(supportType, proxyType) {
		return nil, nil, helpMessage, fmt.Errorf("invalid proxy type: %s, support types: %v", proxyType, supportType)
	}
	pc := v1.NewProxyConfigurerByType(v1.ProxyType(proxyType))
	if pc == nil {
		return nil, nil, helpMessage, fmt.Errorf("new proxy configurer error")
	}
	config.RegisterProxyFlags(cmd, pc, config.WithSSHMode())

	clientCfg := v1.ClientCommonConfig{}
	config.RegisterClientCommonConfigFlags(cmd, &clientCfg, config.WithSSHMode())

	cmd.InitDefaultHelpCmd()
	if err := cmd.ParseFlags(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			helpMessage = cmd.UsageString()
		}
		return nil, nil, helpMessage, err
	}
	// if name is not set, generate a random one
	if pc.GetBaseConfig().Name == "" {
		id, err := util.RandIDWithLen(8)
		if err != nil {
			return nil, nil, helpMessage, fmt.Errorf("generate random id error: %v", err)
		}
		pc.GetBaseConfig().Name = fmt.Sprintf("sshtunnel-%s-%s", proxyType, id)
	}
	return &clientCfg, pc, helpMessage, nil
}

func (s *TunnelServer) handleNewChannel(channel ssh.NewChannel, extraPayloadCh chan string) {
	ch, reqs, err := channel.Accept()
	if err != nil {
		return
	}
	if s.firstChannel == nil {
		s.firstChannel = ch
	}
	go s.keepAlive(ch)

	for req := range reqs {
		if req.WantReply {
			_ = req.Reply(true, nil)
		}
		if req.Type != "exec" || len(req.Payload) <= 4 {
			continue
		}
		end := 4 + binary.BigEndian.Uint32(req.Payload[:4])
		if len(req.Payload) < int(end) {
			continue
		}
		extraPayload := string(req.Payload[4:end])
		// todo 学习
		select {
		case extraPayloadCh <- extraPayload:
		default:
		}
	}
}

func (s *TunnelServer) keepAlive(ch ssh.Channel) {
	tk := time.NewTicker(time.Second * 30)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			_, err := ch.SendRequest("heartbeat", false, nil)
			if err != nil {
				return
			}
		case <-s.doneCh:
			return
		}
	}
}
