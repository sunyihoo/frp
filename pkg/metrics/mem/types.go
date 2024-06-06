package mem

import (
	"github.com/sunyihoo/frp/pkg/util/metric"
	"time"
)

type ProxyStatistics struct {
	Name          string
	ProxyType     string
	TrafficIn     metric.DateCounter
	TrafficOut    metric.DateCounter
	CurConns      metric.Counter
	LastStartTime time.Time
	LastCloseTime time.Time
}

type ServerStatics struct {
	TotalTrafficIn  metric.DateCounter
	TotalTrafficOut metric.DateCounter
	CurConns        metric.Counter

	// 客户计数器 counter for clients
	ClientCounts metric.Counter

	// 代理类型的计数器 counter for proxy types
	ProxyTypeCounts map[string]metric.Counter

	// 不同代理的统计信息
	// key 键名是代理名称
	ProxyStatistics map[string]*ProxyStatistics
}
