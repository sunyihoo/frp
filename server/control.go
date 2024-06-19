package server

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/sunyihoo/frp/pkg/auth"
	"github.com/sunyihoo/frp/pkg/config"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	pkgerr "github.com/sunyihoo/frp/pkg/errors"
	"github.com/sunyihoo/frp/pkg/msg"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/transport"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"github.com/sunyihoo/frp/pkg/util/util"
	"github.com/sunyihoo/frp/pkg/util/xlog"
	"github.com/sunyihoo/frp/server/controller"
	"github.com/sunyihoo/frp/server/metrics"
	"github.com/sunyihoo/frp/server/proxy"
	"net"
	"runtime/debug"
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
	ctl.registerMsgHandlers()
	ctl.msgTransporter = transport.NewMessageTransporter(ctl.msgDispatcher.SendChannel())
	return ctl, nil
}

// GetWorkConn
// 当 frps 获得一个用户连接时，我们从池中获取一个工作连接并返回它。
// 如果池中没有可用的 workConn，请向 frpc 发送消息以获取一个或多个，并等待它可用。
// 如果等待超时，则返回错误
func (ctl *Control) GetWorkConn() (workConn net.Conn, err error) {
	xl := ctl.xl
	defer func() {
		if err := recover(); err != nil {
			xl.Errorf("panic error: %v", err)
			xl.Errorf(string(debug.Stack()))
		}
	}()

	var ok bool
	// get a work connection from the pool
	select {
	case workConn, ok = <-ctl.workConnCh:
		if !ok {
			err = pkgerr.ErrCtlClosed
			return
		}
		xl.Debugf("get work connection from pool")
	default:
		// 池中没有可用的工作连接，请向 frpc 发送消息以获取更多信息
		if err := ctl.msgDispatcher.Send(&msg.ReqWorkConn{}); err != nil {
			return nil, fmt.Errorf("control is already closed")
		}

		select {
		case workConn, ok = <-ctl.workConnCh:
			if !ok {
				err = pkgerr.ErrCtlClosed
				xl.Warnf("no work connections avalible, %v", err)
				return
			}
		case <-time.After(time.Duration(ctl.serverCfg.UserConnTimeout) * time.Second):
			err = fmt.Errorf("timeout trying to get work connection")
			xl.Warnf("%v", err)
			return
		}
	}

	// 当我们从池中获取工作连接时，请将其替换为新的连接。
	_ = ctl.msgDispatcher.Send(&msg.ReqWorkConn{})
	return
}

func (ctl *Control) registerMsgHandlers() {
	ctl.msgDispatcher.RegisterHandler(&msg.NewProxy{}, ctl.handleNewProxy)
	ctl.msgDispatcher.RegisterHandler(&msg.Ping{}, ctl.handlePing)
	ctl.msgDispatcher.RegisterHandler(&msg.NatHoleVisitor{}, msg.AsyncHandler(ctl.handleNatHoleVisitor))
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
		remoteAddr, err = ctl.RegisterProxy(inMsg)
	}

	// register proxy in this control
	resp := &msg.NewProxyResp{
		ProxyName: inMsg.ProxyName,
	}
	if err != nil {
		xl.Warnf("new proxy [%s] type [%s] error: %v", inMsg.ProxyName, inMsg.ProxyType, err)
		resp.Error = util.GenerateResponseErrorString(fmt.Sprintf("new proxy [%s] error", inMsg.ProxyName),
			err, lo.FromPtr(ctl.serverCfg.DetailedErrorsToClient))
	} else {
		resp.RemoteAddr = remoteAddr
		xl.Infof("new proxy [%s] type [%s] success", inMsg.ProxyName, inMsg.ProxyType)
		metrics.Server.NewProxy(inMsg.ProxyName, inMsg.ProxyType)
	}
	_ = ctl.msgDispatcher.Send(resp)
}

func (ctl *Control) handlePing(m msg.Message) {
	xl := ctl.xl
	inMsg := m.(*msg.Ping)

	content := &plugin.PingContent{
		User: plugin.UserInfo{
			User:  ctl.loginMsg.User,
			Metas: ctl.loginMsg.Metas,
			RunID: ctl.loginMsg.RunID,
		},
		Ping: *inMsg,
	}
	retContent, err := ctl.pluginManager.Ping(content)
	if err == nil {
		inMsg = &retContent.Ping
		err = ctl.authVerifier.VerifyPing(inMsg)
	}
	if err != nil {
		xl.Warnf("received invalid ping: %v", err)
		_ = ctl.msgDispatcher.Send(&msg.Pong{
			Error: util.GenerateResponseErrorString("invalid ping", err, lo.FromPtr(ctl.serverCfg.DetailedErrorsToClient)),
		})
		return
	}
	ctl.lastPing.Store(time.Now())
	xl.Debugf("receive heartbeat")
	_ = ctl.msgDispatcher.Send(&msg.Pong{})
}

func (ctl *Control) handleNatHoleVisitor(m msg.Message) {
	inMsg := m.(*msg.NatHoleVisitor)
	ctl.rc.NatHoleController
}

func (ctl *Control) RegisterProxy(pxyMsg *msg.NewProxy) (remoteAddr string, err error) {
	var pxyConf v1.ProxyConfigurer
	// 从 NewProxy 消息加载配置并验证。
	pxyConf, err = config.NewProxyConfigurerFromMsg(pxyMsg, ctl.serverCfg)
	if err != nil {
		return
	}

	// User info
	userInfo := plugin.UserInfo{
		User:  ctl.loginMsg.User,
		Metas: ctl.loginMsg.Metas,
		RunID: ctl.loginMsg.RunID,
	}

	// NewProxy will return an interface Proxy. NewProxy 将返回一个接口 Proxy。
	// 实际上，它根据代理类型创建不同的代理。我们在这里调用 run（）。
	pxy, err := proxy.NewProxy(ctl.ctx, &proxy.Options{
		UserInfo:           userInfo,
		LoginMsg:           ctl.loginMsg,
		PoolCount:          ctl.poolCount,
		ResourceController: ctl.rc,
		GetWorkConnFn:      ctl.GetWorkConn,
		Configurer:         pxyConf,
		ServerCfg:          ctl.serverCfg,
	})
	if err != nil {
		return remoteAddr, err
	}

	// 检查每个客户端中使用的端口数
	if ctl.serverCfg.MaxPortsPerClient > 0 {
		ctl.mu.Lock()
		if ctl.portsUsedNum+pxy.GetUsedPortsNum() > int(ctl.serverCfg.MaxPortsPerClient) {
			ctl.mu.Unlock()
			err = fmt.Errorf("exceed the max_ports_per_client")
			return
		}
		ctl.portsUsedNum += pxy.GetUsedPortsNum()
		ctl.mu.Unlock()

		defer func() {
			if err != nil {
				ctl.mu.Lock()
				ctl.portsUsedNum -= pxy.GetUsedPortsNum()
				ctl.mu.Unlock()
			}
		}()
	}

	if ctl.pxyManager.Exist(pxyMsg.ProxyName) {
		err = fmt.Errorf("proxy [%s] already exists", pxyMsg.ProxyName)
		return
	}

	remoteAddr, err = pxy.Run()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			pxy.Close()
		}
	}()

	err = ctl.pxyManager.Add(pxyMsg.ProxyName, pxy)
	if err != nil {
		return
	}

	ctl.mu.Lock()
	ctl.proxies[pxy.GetName()] = pxy
	ctl.mu.Unlock()
	return
}
