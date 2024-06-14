// Copyright 2020 guylewin, guy@lewin.co.il
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

package tcpmux

import (
	"github.com/sunyihoo/frp/pkg/util/vhost"
	"net"
	"time"
)

type HTTPConnectTCPMuxer struct {
	*vhost.Muxer

	// 如果 passthrough 为 true，则 CONNECT 请求将转发到后端服务。
	// 否则，它将向客户端返回 OK 响应，并将剩余内容转发到后端服务。
	passthrough bool
}

func NewHTTPConnectTcpMuxer(listener net.Listener, passthrough bool, timeout time.Duration) (*HTTPConnectTCPMuxer, error) {
	ret := &HTTPConnectTCPMuxer{passthrough: passthrough}
	mux, err := vhost.NewMuxer(listener, ret.get)
	mux.SetFailHookFunc(ret.auth)
}
