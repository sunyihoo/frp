package visitor

import (
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"sync"
)

type listenerBundle struct {
	l          *netpkg.InternalListener
	sk         string
	allowUsers []string
}

type Manager struct {
	listeners map[string]*listenerBundle

	mu sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		listeners: make(map[string]*listenerBundle),
	}
}
