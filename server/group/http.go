package group

import (
	"github.com/sunyihoo/frp/pkg/util/vhost"
	"sync"
)

type HTTPGroupController struct {
	// 按组名称编制索引的 groups
	groups map[string]*HTTPGroup

	// 将每个组的 createConn 注册到 vhostRouter。
	// createConn 将从组的一个代理获取连接
	vhostRouter *vhost.Routers

	mu sync.Mutex
}

type HTTPGroup struct {
	group           string
	groupKey        string
	domain          string
	location        string
	routeByHTTPUser string

	// createFuncs 按代理名称编制索引的创建方法
	createFuncs map[string]vhost.CreateConnFunc
	pxyNames    []string
	index       uint64
	ctl         *HTTPGroupController
	mu          sync.RWMutex
}
