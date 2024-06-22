package client

import (
	"context"
	"github.com/sunyihoo/frp/pkg/auth"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
	"github.com/sunyihoo/frp/pkg/transport"
	"github.com/sunyihoo/frp/pkg/util/xlog"
	"github.com/sunyihoo/frp/server/proxy"
	"github.com/sunyihoo/frp/server/visitor"
	"net"
	"sync/atomic"
)

type SessionContext struct {
	// The client common configuration
	Common *v1.ClientCommonConfig

	// Unique ID obtained from frps.
	// It should be attached to the login message when reconnecting.
	RunID string
	// 基础控件连接。关闭 conn 后，msgDispatcher 和整个 Control 将退出。
	Conn net.Conn
	// 指示连接是否加密
	ConnEncrypted bool
	// 根据所选方法设置身份验证
	AuthSetter auth.Setter
	// Connector 用于创建新连接，这些连接可以是真正的 TCP 连接或虚拟流。
	Connector Connector
}

type Control struct {
	// service context
	ctx context.Context
	xl  *xlog.Logger

	// session context
	sessionCtx *SessionContext

	// manage all proxies
	pm *proxy.Manager

	// manage all visitors
	vm *visitor.Manager

	doneCh chan struct{}

	// of time.Time, last time got the Pong message
	lastPong atomic.Value

	// msgTransporter 的作用类似于 HTTP2。
	// 它允许在同一控制连接上同时发送多条消息。
	// 服务器的响应消息会根据 laneKey 和消息类型调度到相应的等待 goroutine。
	msgTransporter transport.MessageTransporter

	// msgDispatcher 是控件连接的包装器。
	// 它提供了一个用于发送消息的通道，你可以注册处理程序以根据消息的各自类型处理消息。
	msgDispatcher *msg.Dispatcher
}
