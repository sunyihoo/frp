// Copyright 2023 The frp Authors
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

package v1

import (
	"github.com/samber/lo"
	"os"

	"github.com/sunyihoo/frp/pkg/util/util"
)

type ClientConfig struct {
	ClientCommonConfig

	Proxies  []TypedProxyConfig  `json:"proxies,omitempty"`
	Visitors []TypeVisitorConfig `json:"visitors,omitempty"`
}

type ClientCommonConfig struct {
	APIMetadata

	Auth AuthClientConfig `json:"auth,omitempty"`
	// User 为代理名称指定前缀，以将其与其他客户端区分开来。
	// 如果此值不是“”，则代理名称将自动更改为
	// “{user}”。{proxy_name}“。
	User string `json:"user,omitempty"`

	// ServerAddr 指定要连接到的服务器的地址。
	// 默认情况下，此值为“0.0.0.0”。
	ServerAddr string `json:"serverAddr,omitempty"`
	// ServerPort 指定要连接到服务器的端口。
	// 默认情况下，此值为 7000。
	ServerPort int `json:"serverPort,omitempty"`
	// STUN server，帮助穿透NAT漏洞。
	NatHoleSTUNServer string `json:"natHoleSTUNServer,omitempty"`
	// DNSServer 指定供 FRPC 使用的 DNS 服务器地址。
	// 如果此值为 “”，则将使用默认 DNS。
	DNSServer string `json:"dnsServer,omitempty"`
	// LoginFailExit 控制客户端在登录尝试失败后是否应退出。
	// 如果为 false，客户端将重试，直到登录尝试成功。
	// 默认情况下，此值为 true。
	LoginFailExit *bool `json:"loginFailExit,omitempty"`
	// Start 按名称指定一组已启用的代理。
	// 如果此切片为空，则启用所有提供的代理。
	// 默认情况下，此切片为空。
	Start []string `json:"start,omitempty"`

	Log       LogConfig             `json:"log,omitempty"`
	WebServer WebServerConfig       `json:"webServer,omitempty"`
	Transport ClientTransportConfig `json:"transport,omitempty"`

	// UDPPacketSize 指定 udp 数据包大小。
	// 默认情况下，此值为 1500
	UDPPacketSize int64 `json:"udpPacketSize,omitempty"`
	// 客户端元数据信息
	MetaDatas map[string]string `json:"metaDatas,omitempty"`

	// 包括代理的其他配置文件。
	IncludeConfigFiles []string `json:"includeConfigFiles,omitempty"`
}

func (c *ClientCommonConfig) Complete() {
	c.ServerAddr = util.EmptyOr(c.ServerAddr, "0.0.0.0")
	c.ServerPort = util.EmptyOr(c.ServerPort, 7000)
	c.LoginFailExit = util.EmptyOr(c.LoginFailExit, lo.ToPtr(true))
	c.NatHoleSTUNServer = util.EmptyOr(c.NatHoleSTUNServer, "stun.easyvoip.com:3478")

	c.Auth.Complete()
	c.Log.Complete()
	c.Transport.Complete()
	c.WebServer.Complete()

	c.UDPPacketSize = util.EmptyOr(c.UDPPacketSize, 1500)
}

type ClientTransportConfig struct {
	// Protocol 指定与服务器交互时要使用的协议。
	// 有效值为 “tcp”、“kcp”、“quic”、“websocket” 和 “wss”。
	// 默认情况下，此值为“tcp”。
	Protocol string `json:"protocol,omitempty"`
	// 拨号到服务器等待连接完成的最长时间。
	DialServerTimeout int64 `json:"dialServerTimeout,omitempty"`
	// DialServerKeepAlive 指定 frpc 和 frps 之间活动网络连接的保持活动探测之间的间隔。
	// 如果为负数，则禁用保持活动状态的探头。
	DialServerKeepAlive int64 `json:"dialServerKeepAlive,omitempty"`
	// ConnectServerLocalIP 指定客户端绑定连接到服务器时的地址。
	// 注意：此值仅在 TCP/Websocket 协议中使用。在 KCP 协议中不受支持。
	ConnectServerLocalIP string `json:"connectServerLocalIP,omitempty"`
	// ProxyURL 指定要通过其连接到服务器的代理地址。
	// 如果此值为 “”，则直接连接到服务器。
	// 默认情况下，此值是从“http_proxy”环境变量中读取的。
	ProxyURL string `json:"proxyURL,omitempty"`
	// PoolCount 指定客户端将预先与服务器建立的连接数。
	PoolCount int `json:"poolCount,omitempty"`
	// TCPMux 切换 TCP 流多路复用。
	// 这允许来自客户端的多个请求共享单个 TCP 连接。
	// 如果此值为 true，则服务器必须启用 TCP 多路复用。
	// 默认情况下，此值为 true。
	TCPMux *bool `json:"tcpMux,omitempty"`
	// TCPMuxKeepaliveInterval 指定 TCP stream multiplier 的保持活动间隔。
	// 如果 TCPMux 为 true，则不需要应用层的心跳，因为它只能依赖于 TCPMux 中的心跳。
	TCPMuxKeepaliveInterval int64 `json:"tcpMuxKeepaliveInterval,omitempty"`
	// QUIC 协议选项。
	QUIC QUICOptions `json:"quic,omitempty"`
	// HeartbeatInterval 指定以什么时间间隔将检测信号发送到服务器（以秒为单位）。
	// 不建议更改此值。
	// 默认情况下，此值为 30。设置负值以禁用它。
	HeartbeatInterval int64 `json:"heartbeatInterval,omitempty"`
	// HeartBeatTimeout 指定连接终止前允许的最大检测信号响应延迟（以秒为单位）。
	// 不建议更改此值。
	// 默认情况下，此值为 90。设置负值以禁用它。
	HeartbeatTimeout int64 `json:"heartbeatTimeout,omitempty"`
	// TLS 指定与服务器连接的 TLS 设置。
	TLS TLSClientConfig `json:"tls,omitempty"`
}

func (c *ClientTransportConfig) Complete() {
	c.Protocol = util.EmptyOr(c.Protocol, "tcp")
	c.DialServerTimeout = util.EmptyOr(c.DialServerTimeout, 10)
	c.DialServerKeepAlive = util.EmptyOr(c.DialServerKeepAlive, 7200)
	c.ProxyURL = util.EmptyOr(c.ProxyURL, os.Getenv("http_proxy"))
	c.PoolCount = util.EmptyOr(c.PoolCount, 1)
	c.TCPMux = util.EmptyOr(c.TCPMux, lo.ToPtr(true))
	c.TCPMuxKeepaliveInterval = util.EmptyOr(c.TCPMuxKeepaliveInterval, 60)
	if lo.FromPtr(c.TCPMux) {
		// 如果启用了 TCPMux，则不需要应用层的心跳，因为我们可以依赖 tcpmux 中的心跳。
		c.HeartbeatInterval = util.EmptyOr(c.HeartbeatInterval, -1)
		c.HeartbeatTimeout = util.EmptyOr(c.HeartbeatTimeout, -1)
	} else {
		c.HeartbeatInterval = util.EmptyOr(c.HeartbeatInterval, 30)
		c.HeartbeatTimeout = util.EmptyOr(c.HeartbeatTimeout, 90)
	}
}

type TLSClientConfig struct {
	// TLSEnable 指定在与服务器通信时是否应使用 TLS。
	// 如果“tls.certFile”和“tls.keyFile”有效，
	// 客户端将加载提供的 tls 配置。
	// 从 v0.50.0 开始，默认值已更改为 true，并且默认启用 tls。
	Enable *bool `json:"enable,omitempty"`
	// 如果 DisableCustomTLSFirstByte 设置为 false，
	// 则 frpc 将在启用 tls 时使用第一个自定义字节与 frps 建立连接。
	// 从 v0.50.0 开始，默认值已更改为 true，并且默认禁用第一个自定义字节。
	DisableCustomTLSFirstByte *bool `json:"disableCustomTLSFirstByte,omitempty"`

	TLSConfig
}

func (c *TLSClientConfig) Complete() {
	c.Enable = util.EmptyOr(c.Enable, lo.ToPtr(true))
	c.DisableCustomTLSFirstByte = util.EmptyOr(c.DisableCustomTLSFirstByte, lo.ToPtr(true))
}

type AuthClientConfig struct {
	// Method 指定用于使用 frps 对 frpc 进行身份验证的身份验证方法。
	// 如果指定了“token” - 令牌将被读入登录消息中。
	// 如果指定了“oidc” - 将使用 OIDC 设置颁发 OIDC（Open ID Connect）令牌。
	// 默认情况下，此值为“token”。
	Method AuthMethod `json:"method,omitempty"`
	// 指定是否在其他作用域中包含身份验证信息。
	// 当前支持的作用域包括：“HeartBeats”、“NewWorkConns”。
	AdditionalScopes []string `json:"additionalScopes,omitempty"`
	// Token 指定用于创建要发送到服务器的密钥的授权令牌。
	// 服务器必须具有匹配的令牌才能成功进行授权。
	// 默认情况下，此值为“”。
	Token string               `json:"token,omitempty"`
	OIDC  AuthOIDCClientConfig `json:"oidc,omitempty"`
}

func (c *AuthClientConfig) Complete() {
	c.Method = util.EmptyOr(c.Method, "token")
}

type AuthOIDCClientConfig struct {
	// ClientID 指定用于在 OIDC 身份验证中获取令牌的客户端 ID。
	ClientID string `json:"clientID,omitempty"`
	// ClientSecret 指定用于在 OIDC 身份验证中获取令牌的客户端密钥。
	ClientSecret string `json:"clientSecret,omitempty"`
	// Audience 指定 OIDC 身份验证中令牌的受众。
	Audience string `json:"audience,omitempty"`
	// Scope 指定 OIDC 身份验证中令牌的范围。
	Scope string `json:"scope,omitempty"`
	// TokenEndpointURL 指定实现 OIDC 令牌终结点的 URL。
	// 它将用于获取 OIDC 令牌。
	TokenEndpointURL string `json:"tokenEndpointURL,omitempty"`
	// AdditionalEndpointParams 指定要发送的其他参数
	// 此字段将传输到 OIDC 令牌生成器中的 map[string][]string。
	AdditionalEndpointParams string `json:"additionalEndpointParams,omitempty"`
}
