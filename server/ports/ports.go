package ports

import (
	"sync"
	"time"
)

type PortCtx struct {
	ProxyName  string
	Port       int
	Closed     bool
	UpdateTime time.Time
}

type Manager struct {
	reservedPorts map[string]*PortCtx
	usedPorts     map[int]*PortCtx
	freePorts     map[int]struct{}

	bindAddr string
	netType  string
	mu       sync.Mutex
}
