package server

import (
	"context"
	"crypto/tls"
	"github.com/fatedier/golib/net/mux"
	quic "github.com/quic-go/quic-go"
	"github.com/sunyihoo/frp/pkg/auth"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	modelmetrics "github.com/sunyihoo/frp/pkg/metrics"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/ssh"
	"github.com/sunyihoo/frp/pkg/transport"
	httppkg "github.com/sunyihoo/frp/pkg/util/http"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"github.com/sunyihoo/frp/pkg/util/vhost"
	"github.com/sunyihoo/frp/server/controller"
	"github.com/sunyihoo/frp/server/ports"
	"github.com/sunyihoo/frp/server/proxy"
	"github.com/sunyihoo/frp/server/visitor"
	"net"
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

	return nil, nil
}
