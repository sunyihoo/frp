package tcpmux

import "github.com/sunyihoo/frp/pkg/util/vhost"

type HTTPConnectTCPMuxer struct {
	*vhost.Muxer

	// 如果 passthrough 为 true，则 CONNECT 请求将转发到后端服务。
	// 否则，它将向客户端返回 OK 响应，并将剩余内容转发到后端服务。
	passthrough bool
}
