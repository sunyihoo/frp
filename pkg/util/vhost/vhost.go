// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vhost

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/fatedier/golib/errors"

	"github.com/sunyihoo/frp/pkg/util/log"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"github.com/sunyihoo/frp/pkg/util/xlog"
)

type RouteInfo string

const (
	RouteInfoKey   RouteInfo = "routeInfo"
	RouteConfigKey RouteInfo = "routeConfig"
)

type RequestRouteInfo struct {
	URL        string
	Host       string
	HTTPUser   string
	RemoteAddr string
	URLHost    string
	Endpoint   string
}

type (
	muxFunc         func(net.Conn) (net.Conn, map[string]string, error)
	authFunc        func(conn net.Conn, username, passwd string, reqInfoMap map[string]string) (bool, error)
	hostRewriteFunc func(net.Conn, string) (net.Conn, error)
	successFunc     func(net.Conn, map[string]string) error
	failHookFunc    func(net.Conn)
)

// Muxer 是用于 https 和 tcpmux 代理的功能组件。
// 它接受连接并从连接数据的开头提取虚拟主机信息。
// 然后，它将连接路由到其相应的侦听器。
type Muxer struct {
	listener net.Listener
	timeout  time.Duration

	vhostFunc      muxFunc
	checkAuth      authFunc
	successHook    successFunc
	failHook       failHookFunc
	rewriteHost    hostRewriteFunc
	registryRouter *Routers
}

func NewMuxer(
	listener net.Listener,
	vhostFunc muxFunc,
	timeout time.Duration,
) (mux *Muxer, err error) {
	mux = &Muxer{
		listener:       listener,
		timeout:        timeout,
		vhostFunc:      vhostFunc,
		registryRouter: NewRouters(),
	}
	go mux.run()
	return mux, nil
}

func (v *Muxer) SetCheckAuthFunc(f authFunc) *Muxer {
	v.checkAuth = f
	return v
}

func (v *Muxer) SetSuccessHookFunc(f successFunc) *Muxer {
	v.successHook = f
	return v
}

func (v *Muxer) SetFailHookFunc(f failHookFunc) *Muxer {
	v.failHook = f
	return v
}

func (v *Muxer) SetRewriteHostFunc(f hostRewriteFunc) *Muxer {
	v.rewriteHost = f
	return v
}

type ChooseEndPointFunc func() (string, error)

type CreateConnFunc func(remoteAddr string) (net.Conn, error)

type CreateConnByEndpointFunc func(endpoint, remoteAddr string) (net.Conn, error)

// RouteConfig 是用于匹配 HTTP 请求的参数
type RouteConfig struct {
	Domain          string
	Location        string
	RewriteHost     string
	Username        string
	Password        string
	Headers         map[string]string
	ResponseHeaders map[string]string
	RouteByHTTPUser string

	CreateConnFn           CreateConnFunc
	ChooseEndpointFn       ChooseEndPointFunc
	CreateConnByEndpointFn CreateConnByEndpointFunc
}

// Listen 监听新的域名，如果 rewriteHost 不为空且 rewriteHost func 不为 nil，
// 则将host header重写为 rewriteHost
func (v *Muxer) Listen(ctx context.Context, cfg *RouteConfig) (l *Listener, err error) {
	l = &Listener{
		name:            cfg.Domain,
		location:        cfg.Location,
		routeByHTTPUser: cfg.RouteByHTTPUser,
		rewriteHost:     cfg.RewriteHost,
		username:        cfg.Username,
		password:        cfg.Password,
		mux:             v,
		accept:          make(chan net.Conn),
		ctx:             ctx,
	}
	err = v.registryRouter.Add(cfg.Domain, cfg.Location, cfg.RouteByHTTPUser, l)
	if err != nil {
		return
	}
	return l, nil
}

func (v *Muxer) getListener(name, path, httpUser string) (*Listener, bool) {
	findRouter := func(inName, inPath, inHTTPUser string) (*Listener, bool) {
		vr, ok := v.registryRouter.Get(inName, inPath, inHTTPUser)
		if ok {
			return vr.payload.(*Listener), true
		}
		//尝试检查是否有一个代理未指定 routerByHTTPUser，这意味着匹配所有。
		vr, ok = v.registryRouter.Get(inName, inPath, "")
		if ok {
			return vr.payload.(*Listener), true
		}
		return nil, false
	}

	// 首先，我们检查完整的主机名
	// 如果不存在,则检查 wildcard_domain，例如 .example.com
	l, ok := findRouter(name, path, httpUser)
	if ok {
		return l, true
	}

	domainSplit := strings.Split(name, ".")
	for {
		if len(domainSplit) < 3 {
			break
		}

		domainSplit[0] = "*"
		name = strings.Join(domainSplit, ".")

		l, ok = findRouter(name, path, httpUser)
		if ok {
			return l, true
		}
		domainSplit = domainSplit[1:]
	}
	// 最后，尝试检查是否有一个代理域是“*”表示匹配所有域。
	l, ok = findRouter("*", path, httpUser)
	if ok {
		return l, true
	}
	return nil, false
}

func (v *Muxer) run() {
	for {
		conn, err := v.listener.Accept()
		if err != nil {
			return
		}
		go v.handle(conn)
	}
}

func (v *Muxer) handle(c net.Conn) {
	if err := c.SetDeadline(time.Now().Add(v.timeout)); err != nil {
		_ = c.Close()
		return
	}

	sConn, reqInfoMap, err := v.vhostFunc(c)
	if err != nil {
		log.Debugf("get hostname from http/https request error: %v", err)
		_ = c.Close()
		return
	}

	name := strings.ToLower(reqInfoMap["Host"])
	path := strings.ToLower(reqInfoMap["Path"])
	httpUser := reqInfoMap["HTTPUser"]
	l, ok := v.getListener(name, path, httpUser)
	if !ok {
		log.Debugf("http request for host [%s] path [%s] httpUser [%s] not found", name, path, httpUser)
		v.failHook(sConn)
		return
	}

	xl := xlog.FromContextSafe(l.ctx)
	if v.successHook != nil {
		if err := v.successHook(c, reqInfoMap); err != nil {
			xl.Infof("success func failure on vhost connection: %v", err)
			_ = c.Close()
			return
		}
	}

	// 如果 checkAuth func 存在并且设置了 username/password，
	// 则验证用户访问权限
	if l.mux.checkAuth != nil && l.username != "" {
		ok, err := l.mux.checkAuth(c, l.username, l.password, reqInfoMap)
		if !ok || err != nil {
			xl.Debugf("auth failed for user: %s", l.username)
			_ = c.Close()
			return
		}
	}

	if err = sConn.SetDeadline(time.Time{}); err != nil {
		_ = c.Close()
		return
	}
	c = sConn

	xl.Debugf("new request host [%s] path [%s] httpUser [%s]", name, path, httpUser)
	// todo PanicToError
	err = errors.PanicToError(func() {
		l.accept <- c
	})
	if err != nil {
		xl.Warnf("listener is already closed, ignore this request")
	}
}

type Listener struct {
	name            string
	location        string
	routeByHTTPUser string
	rewriteHost     string
	username        string
	password        string
	mux             *Muxer // for closing Muxer
	accept          chan net.Conn
	ctx             context.Context
}

func (l *Listener) Accept() (net.Conn, error) {
	xl := xlog.FromContextSafe(l.ctx)
	conn, ok := <-l.accept
	if !ok {
		return nil, fmt.Errorf("Listener closed")
	}

	// 如果 rewriteHost func 存在，
	// 则使用修改后的主机标头重写 http 请求
	// 如果 l.rewriteHost 为空，则无所事事
	if l.mux.rewriteHost != nil {
		sConn, err := l.mux.rewriteHost(conn, l.rewriteHost)
		if err != nil {
			xl.Warnf("host header rewrite failed: %v", err)
			return nil, fmt.Errorf("host header rewrite failed")
		}
		xl.Debugf("rewrite host to [%s] success", l.rewriteHost)
		conn = sConn
	}
	return netpkg.NewContextConn(l.ctx, conn), nil
}

func (l *Listener) Close() error {
	l.mux.registryRouter.Del(l.name, l.location, l.routeByHTTPUser)
	close(l.accept)
	return nil
}

func (l *Listener) Name() string {
	return l.name
}

func (l *Listener) Addr() net.Addr {
	return (*net.TCPAddr)(nil)
}
