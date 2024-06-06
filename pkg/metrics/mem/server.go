package mem

import (
	server "github.com/sunyihoo/frp/server/metrics"
	"sync"
)

var (
	sm = new
	ServerMetrics server.ServerMetrics
)

type serverMetrics struct {
	info *ServerStatics
	mu   sync.Mutex
}

func newServerMetrics() *serverMetrics {
	return &serverMetrics{
		info: &ServerStatics{
			TotalTrafficIn: metric.New
		},
	}
}
