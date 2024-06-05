package server

import (
	"context"
	"crypto/tls"
	"github.com/fatedier/golib/net/mux"
	quic "github.com/quic-go/quic-go"
	"github.com/sunyihoo/frp/pkg/auth"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/ssh"
	httppkg "github.com/sunyihoo/frp/pkg/util/http"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"github.com/sunyihoo/frp/pkg/util/vhost"
	"github.com/sunyihoo/frp/server/controller"
	"github.com/sunyihoo/frp/server/proxy"
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
	pluginManger *plugin.Manager

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
	return nil, nil
}
