package proxy

import (
	"context"
	"fmt"
	"github.com/sunyihoo/frp/pkg/config/types"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/util/xlog"
	"github.com/sunyihoo/frp/server/controller"
	"golang.org/x/time/rate"
	"net"
	"reflect"
	"sync"
)

var proxyFactoryRegistry = map[reflect.Type]func(*BaseProxy) Proxy{}

type GetWorkConnFn func() (net.Conn, error)

type Proxy interface {
	Context() context.Context
	Run() (remoteAddr string, err error)
	GetName() string
	GetConfigurer() v1.ProxyConfigurer
	GetWorkConnFromPool(src, dst net.Addr) (workConn net.Conn, err error)
	GetUsedPortsNum() int
	GetResourceController() *controller.ResourceController
	GetUserInfo() plugin.UserInfo
	GetLimiter() *rate.Limiter
	GetLoginMsg() *msg.Login
	Close()
}

type BaseProxy struct {
	name          string
	rc            *controller.ResourceController
	listeners     []net.Listener
	usedPortsNum  int
	poolCount     int
	getWorkConnFn GetWorkConnFn
	serverCfg     *v1.ServerConfig
	limiter       *rate.Limiter
	userInfo      plugin.UserInfo
	loginMsg      *msg.Login
	configurer    v1.ProxyConfigurer

	mu  sync.RWMutex
	xl  *xlog.Logger
	ctx context.Context
}

func (pxy *BaseProxy) GetName() string {
	return pxy.name
}

func (pxy *BaseProxy) Context() context.Context {
	return pxy.ctx
}

func (pxy *BaseProxy) GetUsedPortsNum() int {
	return pxy.usedPortsNum
}

func (pxy *BaseProxy) GetResourceController() *controller.ResourceController {
	return pxy.rc
}

func (pxy *BaseProxy) GetUserInfo() plugin.UserInfo {
	return pxy.userInfo
}

func (pxy *BaseProxy) GetLoginMsg() *msg.Login {
	return pxy.loginMsg
}

func (pxy *BaseProxy) GetConfigurer() v1.ProxyConfigurer {
	return pxy.configurer
}

func (pxy *BaseProxy) Close() {
	xl := xlog.FromContextSafe(pxy.ctx)
	xl.Infof("proxy closing")
	for _, l := range pxy.listeners {
		l.Close()
	}
}

type Options struct {
	UserInfo           plugin.UserInfo
	LoginMsg           *msg.Login
	PoolCount          int
	ResourceController *controller.ResourceController
	GetWorkConnFn      GetWorkConnFn
	Configurer         v1.ProxyConfigurer
	ServerCfg          *v1.ServerConfig
}

func NewProxy(ctx context.Context, options *Options) (pxy Proxy, err error) {
	configurer := options.Configurer
	xl := xlog.FromContextSafe(ctx).AppendPrefix(configurer.GetBaseConfig().Name)

	// todo 学习 rate.Limiter
	var limiter *rate.Limiter
	limitBytes := configurer.GetBaseConfig().Transport.BandwidthLimit.Bytes()
	if limitBytes > 0 && configurer.GetBaseConfig().Transport.BandwidthLimitMode == types.BandwidthLimitModeServer {
		limiter = rate.NewLimiter(rate.Limit(float64(limitBytes)), int(limitBytes))
	}

	basePxy := BaseProxy{
		name:          configurer.GetBaseConfig().Name,
		rc:            options.ResourceController,
		listeners:     make([]net.Listener, 0),
		poolCount:     options.PoolCount,
		getWorkConnFn: options.GetWorkConnFn,
		serverCfg:     options.ServerCfg,
		limiter:       limiter,
		xl:            xl,
		ctx:           xlog.NewContext(ctx, xl),
		userInfo:      options.UserInfo,
		loginMsg:      options.LoginMsg,
		configurer:    configurer,
	}

	factory := proxyFactoryRegistry[reflect.TypeOf(configurer)]
	if factory == nil {
		return pxy, fmt.Errorf("proxy type not support")
	}
	pxy = factory(&basePxy)
	if pxy == nil {
		return nil, fmt.Errorf("proxy not created")
	}
	return pxy, nil
}

type Manager struct {
	// 按代理名称索引的代理
	pxys map[string]Proxy

	mu sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		pxys: make(map[string]Proxy),
	}
}

func (pm *Manager) Add(name string, pxy Proxy) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, ok := pm.pxys[name]; ok {
		return fmt.Errorf("proxy name [%s] is already in use", name)
	}

	pm.pxys[name] = pxy
	return nil
}

func (pm *Manager) Exist(name string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	_, ok := pm.pxys[name]
	return ok
}

func (pm *Manager) Del(name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.pxys, name)
}

func (pm *Manager) GetByName(name string) (pxy Proxy, ok bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pxy, ok = pm.pxys[name]
	return
}
