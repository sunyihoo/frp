package server

import (
	"context"
	"github.com/sunyihoo/frp/pkg/auth"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/transport"
	"github.com/sunyihoo/frp/pkg/util/xlog"
	"github.com/sunyihoo/frp/server/controller"
	"github.com/sunyihoo/frp/server/proxy"
	"net"
	"sync"
	"sync/atomic"
)

type ControlManager struct {
	// 按运行 ID 编制索引的控件(controls)
	ctlsByRunID map[string]*Control

	mu sync.RWMutex
}

func NewControlManager() *ControlManager {
	return &ControlManager{
		ctlsByRunID: make(map[string]*Control),
	}
}

type Control struct {
	// 所有资源管理器和控制器
	rc *controller.ResourceController

	// 代理管理器
	pxyManager *proxy.Manager

	// 插件管理器
	pluginManager *plugin.Manager

	// 根据所选方法验证身份验证
	authVerify auth.Verifier

	// 其他组件可以使用它来与客户端通信
	msgTransporter transport.MessageTransporter

	// msgDispatcher 是控件(control connection)连接的包装器。
	// 它提供了一个用于发送消息的通道，您可以注册处理程序以根据消息的各自类型处理消息。
	msgDispatcher *msg.Dispatcher

	// 登录消息
	loginMsg *msg.Login

	// 控制连接(control connections)
	conn net.Conn

	// 工作连接
	workConnCh chan net.Conn

	// pool count
	poolCount int

	// 使用的端口，用于限制
	portsUsedNum int

	// 上次收到 Ping 消息
	lastPing atomic.Value

	// 当新的客户端登录时，将生成一个新的运行ID。
	// 如果从登录消息获取的运行ID 具有相同的运行ID，则表示它是相同的客户端，
	// 因此我们可以立即替换旧控制器。
	runID string

	mu sync.RWMutex

	// 服务端配置信息
	serverCfg *v1.ServerConfig

	xl     *xlog.Logger
	ctx    context.Context
	doneCh chan struct{}
}
