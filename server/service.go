package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/fatedier/golib/net/mux"
	fmux "github.com/hashicorp/yamux"
	quic "github.com/quic-go/quic-go"
	"github.com/samber/lo"
	"github.com/sunyihoo/frp/pkg/auth"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	modelmetrics "github.com/sunyihoo/frp/pkg/metrics"
	"github.com/sunyihoo/frp/pkg/msg"
	"github.com/sunyihoo/frp/pkg/nathole"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/ssh"
	"github.com/sunyihoo/frp/pkg/transport"
	httppkg "github.com/sunyihoo/frp/pkg/util/http"
	"github.com/sunyihoo/frp/pkg/util/log"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"github.com/sunyihoo/frp/pkg/util/tcpmux"
	"github.com/sunyihoo/frp/pkg/util/util"
	"github.com/sunyihoo/frp/pkg/util/vhost"
	"github.com/sunyihoo/frp/pkg/util/xlog"
	"github.com/sunyihoo/frp/server/controller"
	"github.com/sunyihoo/frp/server/group"
	"github.com/sunyihoo/frp/server/ports"
	"github.com/sunyihoo/frp/server/proxy"
	"github.com/sunyihoo/frp/server/visitor"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	connReadTimeout       time.Duration = 10 * time.Second
	vhostReadWriteTimeout time.Duration = 30 * time.Second
)

type Service struct {
	// 将连接分派到不同的处理程序侦听同一端口。
	muxer *mux.Mux

	// 接受来自客户端的连接。
	listener net.Listener

	// 使用 kcp 接受连接。
	kcpListener net.Listener

	// 使用 quic 接受连接。
	quicListener *quic.Listener

	// 使用 websocket 接受连接。
	websocketListener net.Listener

	// 接受 frp tls 连接
	tlsListener net.Listener

	// 接受来自 ssh 隧道网关的管道连接
	sshTunnelListener *netpkg.InternalListener

	// 管理所有控制器
	ctlManager *ControlManager

	// 管理所有代理
	pxyManager *proxy.Manager

	// 管理所有代理
	pluginManager *plugin.Manager

	// HTTP 虚拟主机路由器
	httpVhostRouter *vhost.Routers

	// 所有资源管理器和控制器
	rc *controller.ResourceController

	// 仪表板 UI 和 API 的 Web 服务器
	webServer *httppkg.Server

	sshTunnelGateWay *ssh.Gateway

	// 根据所选方法验证身份验证
	authVerifier auth.Verifier

	tlsConfig *tls.Config

	cfg *v1.ServerConfig

	// 服务上下文
	ctx context.Context

	// 调用 cancel 以停止服务
	cancel context.CancelFunc
}

func NewService(cfg *v1.ServerConfig) (*Service, error) {
	tlsConfig, err := transport.NewServerTLSConfig(
		cfg.Transport.TLS.CertFile,
		cfg.Transport.TLS.KeyFile,
		cfg.Transport.TLS.TrustedCaFile,
	)
	if err != nil {
		return nil, err
	}

	var webServer *httppkg.Server
	if cfg.WebServer.Port > 0 {
		ws, err := httppkg.NewServer(cfg.WebServer)
		if err != nil {
			return nil, err
		}
		webServer = ws

		modelmetrics.EnableMem()
		if cfg.EnablePrometheus {
			modelmetrics.EnablePrometheus()
		}
	}

	svr := &Service{
		ctlManager:    NewControlManager(),
		pxyManager:    proxy.NewManager(),
		pluginManager: plugin.NewManager(),
		rc: &controller.ResourceController{
			VisitorManager: visitor.NewManager(),
			TCPPortManager: ports.NewManager("tcp", cfg.ProxyBindAddr, cfg.AllowPorts),
			UDPPortManager: ports.NewManager("udp", cfg.ProxyBindAddr, cfg.AllowPorts),
		},
		sshTunnelListener: netpkg.NewInternalListener(),
		httpVhostRouter:   vhost.NewRouters(),
		authVerifier:      auth.NewAuthVerifier(cfg.Auth),
		webServer:         webServer,
		tlsConfig:         tlsConfig,
		cfg:               cfg,
		ctx:               context.Background(),
	}
	if webServer != nil {
		webServer.RouteRegister(svr.registerRouterHandles)
	}

	// Create tcpmux httpconnect multiplexer. 创建 tcpmux httpconnect 多路复用器。
	if cfg.TCPMuxHTTPConnectPort > 0 {
		var l net.Listener
		address := net.JoinHostPort(cfg.ProxyBindAddr, strconv.Itoa(cfg.TCPMuxHTTPConnectPort))
		l, err = net.Listen("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("create server listener error, %v", err)
		}

		svr.rc.TCPMuxHTTPConnectMuxer, err = tcpmux.NewHTTPConnectTCPMuxer(l, cfg.TCPMuxPassthrough, vhostReadWriteTimeout)
		if err != nil {
			return nil, fmt.Errorf("create vhost tcpMuxer error, %v", err)
		}
		log.Infof("tcpmux httpconnect multiplexer listen on %s, passthrough: %v", address, cfg.TCPMuxPassthrough)
	}

	// Init all plugins
	for _, p := range cfg.HTTPPlugins {
		svr.pluginManager.Register(plugin.NewHTTPPluginOptions(p))
		log.Infof("plugin [%s] has been registered", p.Name)
	}
	svr.rc.PluginManager = svr.pluginManager

	// Init group controller
	svr.rc.TCPGroupCtl = group.NewTCPGroupCtl(svr.rc.TCPPortManager)

	// Init HTTP group controller
	svr.rc.HTTPGroupCtl = group.NewHTTPGroupController(svr.httpVhostRouter)

	// Init TCO mux group controller
	svr.rc.TCPMuxGroupCtl = group.NewTCPMuxGroupCtl(svr.rc.TCPMuxHTTPConnectMuxer)

	// Init 404 not found page
	vhost.NotFoundPagePath = cfg.Custom404Page

	var (
		httpMuxOn  bool
		httpsMuxOn bool
	)
	if cfg.BindAddr == cfg.ProxyBindAddr {
		if cfg.BindPort == cfg.VhostHTTPPort {
			httpMuxOn = true
		}
		if cfg.BindPort == cfg.VhostHTTPSPort {
			httpsMuxOn = true
		}
	}

	// 侦听接受来自客户端的连接。
	address := net.JoinHostPort(cfg.BindAddr, strconv.Itoa(cfg.BindPort))
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("create server listener error, %v", err)
	}

	svr.muxer = mux.NewMux(ln)
	svr.muxer.SetKeepAlive(time.Duration(cfg.Transport.TCPKeepAlive) * time.Second)
	go func() {
		_ = svr.muxer.Serve()
	}()
	ln = svr.muxer.DefaultListener()

	svr.listener = ln
	log.Infof("frps tcp listen on %s", address)

	// Listen for accepting connections from client using kcp protocol.
	if cfg.KCPBindPort > 0 {
		address := net.JoinHostPort(cfg.BindAddr, strconv.Itoa(cfg.KCPBindPort))
		svr.kcpListener, err = netpkg.ListenKcp(address)
		if err != nil {
			// todo kcp udp address?
			return nil, fmt.Errorf("listen on kcp udp address %s error: %v", address, err)
		}
		log.Infof("frps kcp listen on udp %s", address)
	}

	if cfg.QUICBindPort > 0 {
		address := net.JoinHostPort(cfg.BindAddr, strconv.Itoa(cfg.QUICBindPort))
		quicTLSCfg := tlsConfig.Clone()
		quicTLSCfg.NextProtos = []string{"frp"}
		svr.quicListener, err = quic.ListenAddr(address, quicTLSCfg, &quic.Config{
			MaxIdleTimeout:     time.Duration(cfg.Transport.QUIC.MaxIdleTimeout) * time.Second,
			MaxIncomingStreams: int64(cfg.Transport.QUIC.MaxIncomingStreams),
			KeepAlivePeriod:    time.Duration(cfg.Transport.QUIC.KeepalivePeriod) * time.Second,
		})
		if err != nil {
			return nil, fmt.Errorf("listen on quic udp address %s error: %v", address, err)
		}
		log.Infof("frps quic listen on %s", address)
	}

	if cfg.SSHTunnelGateway.BindPort > 0 {
		sshGateWay, err := ssh.NewGateway(cfg.SSHTunnelGateway, cfg.ProxyBindAddr, svr.sshTunnelListener)
		if err != nil {
			return nil, fmt.Errorf("create ssh gateway error: %v", err)
		}
		svr.sshTunnelGateWay = sshGateWay
		log.Infof("frps sshTunnelGateway listen on port %d", cfg.SSHTunnelGateway.BindPort)
	}

	// Listen for accepting connections from client using websocket protocol.
	websocketPrefix := []byte("GET " + netpkg.FrpWebsocketPath)
	websocketLn := svr.muxer.Listen(0, uint32(len(websocketPrefix)), func(data []byte) bool {
		return bytes.Equal(data, websocketPrefix)
	})
	svr.websocketListener = websocketLn

	// Create http vhost muxer
	if cfg.VhostHTTPPort > 0 {
		rp := vhost.NewHTTPReverseProxy(vhost.HTTPReverseProxyOptions{
			ResponseHeaderTimeoutS: cfg.VhostHTTPTimeout,
		}, svr.httpVhostRouter)
		svr.rc.HTTPReverseProxy = rp

		address := net.JoinHostPort(cfg.ProxyBindAddr, strconv.Itoa(cfg.VhostHTTPPort))
		server := &http.Server{
			Addr:              address,
			Handler:           rp,
			ReadHeaderTimeout: 60 * time.Second,
		}
		var l net.Listener
		if httpMuxOn {
			l = svr.muxer.ListenHTTP(1)
		} else {
			l, err = net.Listen("tcp", address)
			if err != nil {
				return nil, fmt.Errorf("create vhost http listener error: %v", err)
			}
		}
		go func() {
			_ = server.Serve(l)
		}()
		log.Infof("http service listen on %s", address)
	}

	// Create https vhost muxer
	if cfg.VhostHTTPSPort > 0 {
		var l net.Listener
		if httpsMuxOn {
			l = svr.muxer.ListenHTTPS(1)
		} else {
			address := net.JoinHostPort(cfg.ProxyBindAddr, strconv.Itoa(cfg.VhostHTTPSPort))
			l, err = net.Listen("tcp", address)
			if err != nil {
				return nil, fmt.Errorf("create server listener error, %v", err)
			}
			log.Infof("https service listen on %s", address)
		}

		svr.rc.VhostHTTPSMuxer, err = vhost.NewHTTPSMuxer(l, vhostReadWriteTimeout)
		if err != nil {
			return nil, fmt.Errorf("create vhost hppsMuxer error, %v", err)
		}
	}

	// frps tls listener
	svr.tlsListener = svr.muxer.Listen(2, 1, func(data []byte) bool {
		// 仅当虚拟主机 HTTPS 端口与绑定端口不同时，TLS第一个字节才能是0x16
		return int(data[0]) == netpkg.FRPTLSHeadByte || int(data[0]) == 0x16
	})

	// Create nat hole controller
	nc, err := nathole.NewController(time.Duration(cfg.NatHoleAnalysisDataReserveHours) * time.Hour)
	if err != nil {
		return nil, fmt.Errorf("create nat hole controller error, %v", err)
	}
	svr.rc.NatHoleController = nc

	return svr, nil
}

func (svr *Service) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	svr.ctx = ctx
	svr.cancel = cancel

	// run dashboard web server
	if svr.webServer != nil {
		go func() {
			log.Infof("dashboard listen on %s", svr.webServer.Address())
			if err := svr.webServer.Run(); err != nil {
				log.Warnf("dashboard server exot with error: %v", err)
			}
		}()
	}

	go svr.HandleListener(svr.sshTunnelListener, true)
}

func (svr *Service) handleConnection(ctx context.Context, conn net.Conn, internal bool) {
	xl := xlog.FromContextSafe(ctx)

	var (
		rawMsg msg.Message
		err    error
	)

	_ = conn.SetReadDeadline(time.Now().Add(connReadTimeout))
	if rawMsg, err = msg.ReadMsg(conn); err != nil {
		log.Tracef("Failed to read message: %v", err)
		conn.Close()
		return
	}
	_ = conn.SetReadDeadline(time.Time{})

	switch m := rawMsg.(type) {
	case *msg.Login:
		// 服务器插件钩子 server plugin hook
		content := &plugin.LoginContent{
			Login:         *m,
			ClientAddress: conn.RemoteAddr().String(),
		}
		retContent, err := svr.pluginManager.Login(content)
		if err == nil {
			m = &retContent.Login
			err = svr.RegisterControl(conn, m, internal)
		}

	}

}

// HandleListener 接受来自客户端的连接，并调用 handleConnection 来处理它们。
// 如果 internal 为 true，则表示此侦听器用于内部通信，如 ssh 隧道网关。
// TODO(fatedier): Pass some parameters of listener/connection through context to avoid passing too many parameters. TODO（fatedier）：通过上下文传递 listenerconnection 的一些参数，避免传递过多的参数。
func (svr *Service) HandleListener(l net.Listener, internal bool) {
	// Listen for incoming connections from client.
	for {
		c, err := l.Accept()
		if err != nil {
			log.Warnf("Listen for incoming connections from client closed")
			return
		}
		// inject xlog object into net.Conn context
		xl := xlog.New()
		ctx := context.Background()

		c = netpkg.NewContextConn(xlog.NewContext(ctx, xl), c)

		if !internal {
			log.Tracef("start check TLS connection...")
			originConn := c
			forceTLS := svr.cfg.Transport.TLS.Force
			var isTLS, custom bool
			c, isTLS, custom, err = netpkg.CheckAndEnableTLSServerConnWithTimeout(c, svr.tlsConfig, forceTLS, connReadTimeout)
			if err != nil {
				log.Warnf("CheckAndEnableTLSServerConnWithTimeout error: %v", err)
				originConn.Close()
				continue
			}
			log.Tracef("check TLS connection success,isTLS: %v custom: %v internal: %v", isTLS, custom, internal)
		}

		// Start a new goroutine to handle connection.
		go func(ctx context.Context, frpConn net.Conn) {
			if lo.FromPtr(svr.cfg.Transport.TCPMux) && !internal {
				fmuxCfg := fmux.DefaultConfig()
				fmuxCfg.KeepAliveInterval = time.Duration(svr.cfg.Transport.TCPMuxKeepaliveInternal) * time.Second
				fmuxCfg.LogOutput = io.Discard
				fmuxCfg.MaxStreamWindowSize = 6 * 1024 * 1024
				session, err := fmux.Server(frpConn, fmuxCfg)
				if err != nil {
					log.Warnf("Failed to create mux connection: %v", err)
					frpConn.Close()
					return
				}

				for {
					stream, err := session.AcceptStream()
					if err != nil {
						log.Debugf("Accept new mux stream error: %v", err)
						session.Close()
						return
					}
					go svr.handleConnection(ctx, stream, internal)
				}
			} else {
				svr.handleConnection(ctx, frpConn, internal)
			}
		}(ctx, c)
	}
}

func (svr *Service) RegisterControl(ctlConn net.Conn, loginMsg *msg.Login, internal bool) error {
	// 如果客户端的 RunID 为空，则它是一个新客户端，我们只是创建一个新控制器。
	// 否则，我们会检查是否有一个控制器具有相同的运行 ID。如果是这样，我们释放以前的控制器并启动新的控制器。
	var err error
	if loginMsg.RunID == "" {
		loginMsg.RunID, err = util.RandID()
		if err != nil {
			return err
		}
	}

	ctx := netpkg.NewContextFromConn(ctlConn)
	xl := xlog.FromContextSafe(ctx)
	xl.AppendPrefix(loginMsg.RunID)
	ctx = xlog.NewContext(ctx, xl)
	xl.Infof("client login info: ip [%s] version [%s] hostname [%s] os [%s] arch [%s]",
		ctlConn.RemoteAddr().String(), loginMsg.Version, loginMsg.Hostname, loginMsg.Os, loginMsg.Arch,
	)

	// Check auth
	authVerifier := svr.authVerifier
	if internal && loginMsg.ClientSpec.AlwaysAuthPass {
		authVerifier = auth.AlwaysPassVerifier
	}
	if err := authVerifier.VerifyLogin(loginMsg); err != nil {
		return err
	}

	// TODO(fatedier): use SessionContext
	ctl, err := NewControl()
}
