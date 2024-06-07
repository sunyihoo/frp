package ports

import (
	"github.com/sunyihoo/frp/pkg/config/types"
	"sync"
	"time"
)

const (
	MinPort                    = 1
	MaxPort                    = 65535
	MaxPortReservedDuration    = time.Duration(24) * time.Hour
	CleanReservedPortsInterval = time.Hour
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

func NewManager(netType string, bindAddr string, allowPorts []types.PortsRange) *Manager {
	pm := &Manager{
		reservedPorts: make(map[string]*PortCtx),
		usedPorts:     make(map[int]*PortCtx),
		freePorts:     make(map[int]struct{}),
		bindAddr:      bindAddr,
		netType:       netType,
	}
	if len(allowPorts) > 0 {
		for _, pair := range allowPorts {
			if pair.Single > 0 {
				pm.freePorts[pair.Single] = struct{}{}
			} else {
				for i := pair.Start; i < pair.End; i++ {
					pm.freePorts[i] = struct{}{}
				}
			}
		}
	} else {
		for i := MinPort; i < MaxPort; i++ {
			pm.freePorts[i] = struct{}{}
		}
	}
	go pm.cleanReservedPortsWorker()
	return pm
}

// 如果在过去 24 小时内未使用保留端口，释放该端口。
func (pm *Manager) cleanReservedPortsWorker() {
	for {
		time.Sleep(CleanReservedPortsInterval)
		pm.mu.Lock()
		for name, ctx := range pm.reservedPorts {
			if ctx.Closed && time.Since(ctx.UpdateTime) > MaxPortReservedDuration {
				delete(pm.reservedPorts, name)
			}
		}
		pm.mu.Lock()
	}
}
