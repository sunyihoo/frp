package server

import (
	"context"
	"github.com/sunyihoo/frp/pkg/auth"
	"github.com/sunyihoo/frp/pkg/config"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/transport"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"github.com/sunyihoo/frp/pkg/util/xlog"
	"github.com/sunyihoo/frp/server/controller"
	"github.com/sunyihoo/frp/server/proxy"
	"net"
	"sync"
	"sync/atomic"
	"time"
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
	authVerifier auth.Verifier

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

	// proxies in one client
	proxies map[string]proxy.Proxy

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

// TODO(fatedier): Referencing the implementation of frpc, encapsulate the input parameters as SessionContext.
func NewControl(
	ctx context.Context,
	rc *controller.ResourceController,
	pxyManager *proxy.Manager,
	pluginManager *plugin.Manager,
	authVerifier auth.Verifier,
	ctlConn net.Conn,
	ctlConnEncrypted bool,
	loginMsg *msg.Login,
	serverCfg *v1.ServerConfig,
) (*Control, error) {
	poolCount := loginMsg.PoolCount
	if poolCount > int(serverCfg.Transport.MaxPoolCount) {
		poolCount = int(serverCfg.Transport.MaxPoolCount)
	}
	ctl := &Control{
		rc:            rc,
		pxyManager:    pxyManager,
		pluginManager: pluginManager,
		authVerifier:  authVerifier,
		conn:          ctlConn,
		loginMsg:      loginMsg,
		workConnCh:    make(chan net.Conn, poolCount+10),
		proxies:       make(map[string]proxy.Proxy),
		poolCount:     poolCount,
		portsUsedNum:  0,
		runID:         loginMsg.RunID,
		serverCfg:     serverCfg,
		xl:            xlog.FromContextSafe(ctx),
		ctx:           ctx,
		doneCh:        make(chan struct{}),
	}
	ctl.lastPing.Store(time.Now())

	if ctlConnEncrypted {
		cryptoRW, err := netpkg.NewCryptoReadWriter(ctl.conn, []byte(ctl.serverCfg.Auth.Token))
		if err != nil {
			return nil, err
		}
		ctl.msgDispatcher = msg.NewDispatcher(cryptoRW)
	} else {
		ctl.msgDispatcher = msg.NewDispatcher(ctl.conn)
	}
	ctl.re

}

func (ctl *Control) registerMsgHandlers() {
	ctl.msgDispatcher.RegisterHandler(&msg.NewProxy{}, ctl.handle)
}

func (ctl *Control) handleNewProxy(m msg.Message) {
	xl := ctl.xl
	inMsg := m.(*msg.NewProxy)

	content := &plugin.NewProxyContent{
		User: plugin.UserInfo{
			User:  ctl.loginMsg.User,
			Metas: ctl.loginMsg.Metas,
			RunID: ctl.loginMsg.RunID,
		},
		NewProxy: *inMsg,
	}
	var remoteAddr string
	retContent, err := ctl.pluginManager.NewProxy(content)
	if err == nil {
		inMsg = &retContent.NewProxy
		remoteAddr, err = ctl.Re
	}
}

func (ctl *Control) RegisterProxy(pxyMsg *msg.NewProxy) (remoteAddr string, err error) {
	var pxyConf v1.ProxyConfigurer
	// 从 NewProxy 消息加载配置并验证。
	pxyConf, err = config.New
}
