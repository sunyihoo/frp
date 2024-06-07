// Copyright 2019 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"github.com/sunyihoo/frp/pkg/nathole"
	plugin "github.com/sunyihoo/frp/pkg/plugin/server"
	"github.com/sunyihoo/frp/pkg/util/tcpmux"
	"github.com/sunyihoo/frp/pkg/util/vhost"
	"github.com/sunyihoo/frp/server/group"
	"github.com/sunyihoo/frp/server/ports"
	"github.com/sunyihoo/frp/server/visitor"
)

// ResourceController 所有资源管理器和控制器
type ResourceController struct {
	// 管理所有访客侦听器
	VisitorManager *visitor.Manager

	// TCP 组控制器
	TCPGroupCtl *group.TCPGroupCtl

	// HTTP 组控制器
	HTTPGroupCtl *group.HTTPGroupController

	// TCP 多路复用器组控制器
	TCPMuxGroupCtl *group.TCPMuxGroupCtl

	// 管理所有 TCP 端口
	TCPPortManager *ports.Manager

	// 管理所有 UDP 端口
	UDPPortManager *ports.Manager

	// 对于 HTTP 代理，转发 HTTP 请求
	HTTPReverseProxy *vhost.HTTPReverseProxy

	// 对于 HTTPS 代理，按主机名和其他信息将请求路由到不同的客户端
	VhostHTTPSMuxer *vhost.HTTPMuxer

	// 用于连接nat hole的控制器
	NatHoleController *nathole.Controller

	// 利用 HTTP CONNECT 方法在一个 TCP 连接上多路复用多个流 todo ?
	TCPMuxController *tcpmux.HTTPConnectTCPMuxer

	// 所有服务端管理者插件
	PluginManager *plugin.Manager
}
