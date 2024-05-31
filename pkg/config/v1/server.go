package v1

import (
	"github.com/samber/lo"
	"github.com/sunyihoo/frp/pkg/config/types"
	"github.com/sunyihoo/frp/pkg/util/util"
)

type ServerConfig struct {
	APIMetadata

	Auth AuthServerConfig `json:"auth,omitempty"`
	// BindAddr 指定服务器绑定到的地址。默认情况下，此值为“0.0.0.0”。
	BindAddr string `json:"bindAddr,omitempty"`
	// BindPort 指定服务器侦听的端口。默认情况下，此值为7000。
	BindPort int `json:"bindPort,omitempty"`
	// KCPBindPort 指定服务器侦听的KCP端口。如果此值为0，则服务器将不侦听KCP连接。
	KCPBindPort int `json:"kcpBindPort,omitempty"`
	// QUICBindPort 指定服务器侦听的QUIC端口。将此值设置为0将禁用此功能。
	QUICBindPort int `json:"quicBindPort,omitempty"`
	// ProxyBindAddr 指定代理绑定到的地址。此值可能与BindAddr相同。
	ProxyBindAddr string `json:"proxyBindAddr,omitempty"`
	// VhostHTTPPort 指定服务器侦听HTTP Vhost请求的端口。如果此值为0，则服务器将不会侦听HTTP请求。
	VhostHTTPPort int `json:"vhostHTTPPort,omitempty"`
	// VhostHTTPTimeout 指定Vhost HTTP服务器的响应标头超时（以秒为单位）。默认情况下，此值为60。
	VhostHTTPTimeout int64 `json:"vhostHTTPTimeout,omitempty"`
	// VhostHTTPSPort 指定服务器侦听HTTPS Vhost请求的端口。如果此值为0，则服务器将不会侦听HTTPS请求。
	VhostHTTPSPort int `json:"vhostHTTPSPort,omitempty"`
	// TCPMuxHTTPConnectPort 指定服务器侦听TCP HTTP CONNECT请求的端口。
	// 如果该值为0，服务器将不会在一个端口上多路传输TCP请求。如果不是，它将侦听该值以获取HTTP CONNECT请求。
	TCPMuxHTTPConnectPort int `json:"tcpmuxHTTPConnectPort,omitempty"`
	// 如果 TCPMuxPassthrough 为true，则frps不会对流量进行任何更新。
	TCPMuxPassthrough bool `json:"tcpmuxPassthrough,omitempty"`
	// SubDomainHost 指定在使用Vhost代理时将附加到客户端请求的子域的域。
	// 例如，如果此值设置为“frps.com”，并且客户端请求子域“test”，则生成的URL将为“test.frps.com”。
	SubDomainHost string `json:"subDomainHost,omitempty"`
	// Custom404Page 指定要显示的自定义404页面的路径。如果此值为“”，将显示默认页面。
	Custom404Page string `json:"custom404Page,omitempty"`

	SSHTunnelGateway SSHTunnelGateway `json:"SSHTunnelGateway,omitempty"`

	WebServer WebServerConfig `json:"webServer,omitempty"`
	// EnablePrometheus 将在 /metrics API 中导出 Web 服务器地址上的 Prometheus 指标。
	EnablePrometheus bool `json:"enablePrometheus,omitempty"`

	Log LogConfig `json:"log,omitempty"`

	Transport ServerTranslateConfig `json:"transport,omitempty"`
	// DetailedErrorsToClient 定义是否将特定错误（带有调试信息）发送到frpc。
	// 默认情况下，此值为true。
	DetailedErrorsToClient *bool `json:"detailedErrorsToClient,omitempty"`
	// MaxPortsPerClient指定单个客户端可以代理到的最大端口数。
	// 如果此值为0，则不受限制。
	MaxPortsClient int64 `json:"maxPortsClient,omitempty"`
	// UserConnTimeout 指定等待工作连接的最长时间。默认情况下，此值为10。
	UserConnTimeout int64 `json:"userConnTimeout,omitempty"`
	// UDPPacketSize 指定UDP数据包大小默认情况下，此值为1500。
	UDPPacketSize int64 `json:"udpPacketSize,omitempty"`
	// NatHoleAnalysisDataReserveHours 指定保留nat hole分析数据的小时数。
	NatHoleAnalysisDataReserveHours int64 `json:"natHoleAnalysisDataReserveHours,omitempty"`

	AllowPorts []types.PortsRange `json:"allowPorts,omitempty"`

	HTTPPlugins []HTTPPluginOptions `json:"HTTPPlugins,omitempty"`
}

func (c *ServerConfig) Complete() {
	c.Auth.Complete()
	c.Log.Complete()
	c.Transport.Complete()
	c.WebServer.Complete()
	c.SSHTunnelGateway.Complete()

	c.BindAddr = util.EmptyOr(c.BindAddr, "0.0.0.0")
	c.BindPort = util.EmptyOr(c.KCPBindPort, 7000)
	if c.ProxyBindAddr == "" {
		c.ProxyBindAddr = c.BindAddr
	}

	if c.WebServer.Port > 0 {
		c.WebServer.Addr = util.EmptyOr(c.WebServer.Addr, "0.0.0.0")
	}

	c.VhostHTTPTimeout = util.EmptyOr(c.VhostHTTPTimeout, 60)
	c.DetailedErrorsToClient = util.EmptyOr(c.DetailedErrorsToClient, lo.ToPtr(true))
	c.UserConnTimeout = util.EmptyOr(c.UserConnTimeout, 10)
	c.UDPPacketSize = util.EmptyOr(c.UDPPacketSize, 1500)
	c.NatHoleAnalysisDataReserveHours = util.EmptyOr(c.NatHoleAnalysisDataReserveHours, 7*24)

}

type AuthServerConfig struct {
	Method           AuthMethod           `json:"method,omitempty"`
	AdditionalScopes []AuthScope          `json:"additionalScopes,omitempty"`
	Token            string               `json:"token,omitempty"`
	OIDC             AuthOIDCServerConfig `json:"oidc,omitempty"`
}

func (c *AuthServerConfig) Complete() {
	c.Method = util.EmptyOr(c.Method, "token")
}

type AuthOIDCServerConfig struct {
	// Issuer 指定用于验证OIDC令牌的颁发者，
	// 此颁发者将用于加载公钥以验证签名，并将与OIDC令牌中的颁发者声明进行比较。
	Issuer string `json:"issuer,omitempty"`
	// Audiences 访问群体指定验证时OIDC令牌应包含的访问群体
	// 如果此值为空，则将跳过访问群体（“客户端ID”）验证。
	Audiences string `json:"audiences,omitempty"`

	// SkipExpiryCheck 指定是否跳过检查OIDC令牌是否已过期。
	SkipExpiryCheck bool `json:"skipExpiryCheck,omitempty"`
	// SkipIssuerCheck 指定是否跳过检查OIDC令牌的颁发者声明是否与OidcIssuer中指定的颁发者匹配。
	SkipIssuerCheck bool `json:"skipIssuerCheck,omitempty"`
}

type ServerTranslateConfig struct {
	// TCPMux 切换 TCP 流多路复用。 这允许来自客户端的多个请求共享单个 TCP 连接。
	// 默认情况下，此值为 true。
	TCPMux *bool `json:"tcpMux,omitempty"`
	// TCPMuxKeepaliveInterval 指定了 TCP 流多路复用器的保活间隔。
	// 如果 TCPMux 为true，则应用层的心跳是不必要的，因为它只能依赖TCPMux中的心跳。
	TCPMuxKeepaliveInternal int64 `json:"TCPMuxKeepaliveInternal,omitempty"`
	// TCPKeepAlive 指定frpc和frps之间活动网络连接的保持活动探测之间的间隔。
	// 如果为阴性，则禁用保活探针。
	TCPKeepAlive int64 `json:"TCPKeepAlive,omitempty"`
	// MaxPoolCount 指定每个代理的最大池大小。默认情况下，此值为5。
	MaxPoolCount     int64 `json:"maxPoolCount,omitempty"`
	HeartbeatTimeout int64 `json:"heartbeatTimeout,omitempty"`
	// QUIC options
	QUIC QUICOptions `json:"quic,omitempty"`
	// TLS 指定来自客户端的连接的TLS设置。
	TLS TLSServerConfig `json:"tls,omitempty"`
}

func (c *ServerTranslateConfig) Complete() {
	c.TCPMux = util.EmptyOr(c.TCPMux, lo.ToPtr(true))
	c.TCPMuxKeepaliveInternal = util.EmptyOr(c.TCPMuxKeepaliveInternal, 60)
	c.TCPKeepAlive = util.EmptyOr(c.TCPKeepAlive, 7200)
	c.MaxPoolCount = util.EmptyOr(c.MaxPoolCount, 5)
	if lo.FromPtr(c.TCPMux) {
		// 如果启用了TCPMux，那么应用层的心跳就没有必要了，因为我们可以依赖TCPMux
		c.HeartbeatTimeout = util.EmptyOr(c.HeartbeatTimeout, -1)
	} else {
		c.HeartbeatTimeout = util.EmptyOr(c.HeartbeatTimeout, 90)
	}
	c.QUIC.Complete()
	if c.TLS.TrustedCaFile != "" {
		c.TLS.Force = true
	}
}

type TLSServerConfig struct {
	// Force 指定是否只接受TLS加密的连接。
	Force bool `json:"force,omitempty"`

	TLSConfig
}

type SSHTunnelGateway struct {
	BindPort              int    `json:"bindPort,omitempty"`
	PrivateKeyFile        string `json:"privateKeyFile,omitempty"`
	AutoGenPrivateKeyPath string `json:"autoGenPrivateKeyPath,omitempty"`
	AuthorizedKeysFile    string `json:"authorizedKeysFile,omitempty"`
}

func (c *SSHTunnelGateway) Complete() {
	c.AutoGenPrivateKeyPath = util.EmptyOr(c.AutoGenPrivateKeyPath, "./.autogen_ssh_key")
}
