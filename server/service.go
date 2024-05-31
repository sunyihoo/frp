package server

import (
	"github.com/fatedier/golib/net/mux"
	quic "github.com/quic-go/quic-go"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	netokg "github.com/sunyihoo/frp/pkg/util/net"
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
	sshTunnelListener *netokg.InternalListener

	// 管理所有控制器
	ctlManager *ControlManager

	pxyManager *proxy.Ma
}

func NewService(cfg *v1.ServerConfig) (*Service, error) {
	return nil, nil
}
