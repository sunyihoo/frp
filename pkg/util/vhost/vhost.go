package vhost

import (
	"net"
	"time"
)

type (
	muxFunc         func(net.Conn) (net.Conn, map[string]string, error)
	authFunc        func(conn net.Conn, passwd string, reqInfoMap map[string]string) (bool, error)
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

type ChooseEndPointFunc func() (string, error)

type CreateConnFunc func(remoteAddr string) (net.Conn, error)

type CreateConnByEndpointFunc func(endpoint, remoteAddr string) (net.Conn, error)

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
