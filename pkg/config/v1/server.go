package v1

type ServerConfig struct {
	APIMetadata

	AuthServerConfig `json:"auth,omitempty"`
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
	VhostHTTPTimeout int `json:"vhostHTTPTimeout,omitempty"`
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
}

type AuthServerConfig struct {
	Method           AuthMethod           `json:"method,omitempty"`
	AdditionalScopes []AuthScope          `json:"additionalScopes,omitempty"`
	Token            string               `json:"token,omitempty"`
	OIDC             AuthOIDCServerConfig `json:"oidc,omitempty"`
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
