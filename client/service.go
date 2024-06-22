package client

import (
	"context"
	"github.com/sunyihoo/frp/pkg/auth"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
	httppkg "github.com/sunyihoo/frp/pkg/util/http"
	"net"
	"sync"
	"time"
)

type Service struct {
	ctlMu sync.RWMutex
	// 管理器控制与服务器的连接
	ctl *Control
	//  Uniq id 是从 frps 获取的，它将附加到 loginMsg。
	runID string

	// 根据所选方法设置身份验证
	authSetter auth.Setter

	// web server fpr admin UI and apis
	webServer *httppkg.Server

	cfgMu      sync.RWMutex
	common     *v1.ClientCommonConfig
	proxyCfgs  []v1.ProxyConfigurer
	visitors   []v1.VisitorConfigurer
	clientSpec *msg.ClientSpec

	// 用于初始化此客户端的配置文件，
	// 如果未使用配置文件，则为空字符串。
	configFilePath string

	// service context
	ctx context.Context
	// call cancel to stop service
	cancel                   context.CancelCauseFunc
	gracefulShutdownDuration time.Duration

	connectorCreator func(context.Context, *v1.ClientCommonConfig) Connector
	handleWorkConnCb func(*v1.ProxyBaseConfig, net.Conn, *msg.StartWorkConn) bool
}
