package group

import (
	"github.com/sunyihoo/frp/pkg/util/tcpmux"
	"net"
	"sync"
)

type TCPMuxGroupCtl struct {
	groups map[string]*TCPMuxGroup

	// tcpMuxHTTPConnectMuxer 被用于管理 muxer
	tcpMuxHTTPConnectMuxer *tcpmux.HTTPConnectTCPMuxer
	mu                     sync.Mutex
}

// NewTCPMuxGroupCtl return a new TCPMuxGroupCtl
func NewTCPMuxGroupCtl(tcpMuxHTTPConnectTcpMuxer *tcpmux.HTTPConnectTCPMuxer) *TCPMuxGroupCtl {
	return &TCPMuxGroupCtl{
		groups:                 make(map[string]*TCPMuxGroup),
		tcpMuxHTTPConnectMuxer: tcpMuxHTTPConnectTcpMuxer,
	}
}

type TCPMuxGroup struct {
	group           string
	groupKey        string
	domain          string
	routeByHTTPUser string
	username        string
	password        string

	acceptCh chan net.Conn
	tcpMuxLn net.Listener
	lns      []*TCPMuxGroupListener
	ctl      *TCPMuxGroupCtl
	mu       sync.Mutex
}

// TCPMuxGroupListener TCPMux组侦听者
type TCPMuxGroupListener struct {
	groupName string
	group     *TCPMuxGroup

	addr    net.Addr
	closeCh chan struct{}
}
