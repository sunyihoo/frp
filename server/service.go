package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/fatedier/golib/net/mux"
	quic "github.com/quic-go/quic-go"
	"github.com/sunyihoo/frp/pkg/auth"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	modelmetrics "github.com/sunyihoo/frp/pkg/metrics"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/ssh"
	"github.com/sunyihoo/frp/pkg/transport"
	httppkg "github.com/sunyihoo/frp/pkg/util/http"
	"github.com/sunyihoo/frp/pkg/util/log"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"github.com/sunyihoo/frp/pkg/util/tcpmux"
	"github.com/sunyihoo/frp/pkg/util/vhost"
	"github.com/sunyihoo/frp/server/controller"
	"github.com/sunyihoo/frp/server/group"
	"github.com/sunyihoo/frp/server/ports"
	"github.com/sunyihoo/frp/server/proxy"
	"github.com/sunyihoo/frp/server/visitor"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
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

	sshTunnelGateWay *ssh.GateWay

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

	}

	return nil, nil
}
