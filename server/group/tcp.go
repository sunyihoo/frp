package group

import (
	"github.com/sunyihoo/frp/server/ports"
	"net"
	"sync"
)

type TCPGroupCtl struct {
	groups map[string]*TCPGroup

	// portManager 用于管理端口
	portManager *ports.Manager
	mu          sync.Mutex
}

// TCPGroup 将路由连接到不同的代理
type TCPGroup struct {
	group    string
	groupKey string
	addr     string
	port     string
	realPort string

	acceptCh chan net.Conn
	tcpLn    net.Listener
	lns      []*TCPGroupListener
	ctl      *TCPGroupCtl
	mu       sync.Mutex
}

// TCPGroupListener TCP组侦听者
type TCPGroupListener struct {
	groupName string
	group     *TCPGroup

	addr    net.Addr
	closeCh chan struct{}
}
