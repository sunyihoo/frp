package mem

import (
	"github.com/sunyihoo/frp/pkg/util/metric"
	server "github.com/sunyihoo/frp/server/metrics"
	"sync"
)

var (
	sm            = new
	ServerMetrics server.ServerMetrics
)

type serverMetrics struct {
	info *ServerStatics
	mu   sync.Mutex
}

func newServerMetrics() *serverMetrics {
	return &serverMetrics{
		info: &ServerStatics{
			TotalTrafficIn:  metric.NewDateCounter(ReserveDays),
			TotalTrafficOut: metric.NewDateCounter(ReserveDays),
			CurConns:        metric.NewCounter(),

			ClientCounts:    metric.NewCounter(),
			ProxyTypeCounts: make(map[string]metric.Counter),

			ProxyStatistics: make(map[string]*ProxyStatistics),
		},
	}
}
