package legacy

import legacyauth "github.com/sunyihoo/frp/pkg/auth/legacy"

// ServerCommonConf 包含服务器服务的信息。建议使用 GetDefaultServerConf 而不是直接创建此对象,
// 以便所有未指定的字段都具有合理的默认值。
type ServerCommonConf struct {
	legacyauth.ServerConfig `ini:",extends"`

	// BindAddr 指定服务端绑定的地址。默认情况下，此值为 "0.0.0.0"
	BindAddr string `ini:"bind_addr" json:"bind_addr"`
	// BindPort 指定服务端绑定的端口。默认情况下，此值为 7000
	BindPort int `ini:"bind_port" json:"bind_port"`
	// KCPBindPort 指定服务端监听的KCP端口。如果此值为0，则服务端将不会监听KCP连接。
	// 默认情况下，此值为0。
	KCPBindPort int `ini:"kcp_bind_port" json:"kcp_bind_port"`
	// QUICBindPort 指定服务端监听的QUIC端口。如果此值为0，则服务端将禁用此功能。
	// 默认情况下，此值为0。
	QUICBindPort int `ini:"quci_bind_port" json:"quic_bind_port"`
	// QUICKeepalivePeriod QUIC protocol options
	QUICKeepalivePeriod    int `ini:"quic_keepalive_period" json:"quic_keepalive_period"`
	QUICMaxIdleTimeout     int `ini:"quic_max_idle_timeout" json:"quic_max_idle_timeout"`
	QUICMaxIncomingStreams int `ini:"quic_max_incoming_streams" json:"quic_max_incoming_streams"`
	// ProxyBindAddr 指定代理绑定的地址。默认情况下，此值可能和 BindAddr 相同。
	ProxyBindAddr string `ini:"proxy_bind_addr" json:"proxy_bind_addr"`
	// VhostHTTPPort 指定服务端监听HTTP Vhost请求的端口。如果此值为0，服务端将不会监听HTTP请求。
	// 默认情况下，此值为0。
	VhostHTTPPort int `ini:"vhost_http_port" json:"vhost_http_port"`
	// VhostHTTPSPort 指定服务端监听HTTPS Vhost请求的端口。如果此值为0，服务端将不会监听HTTPS请求。
	// 默认情况下，此值为0。
	VhostHTTPSPort int `ini:"vhost_https_port" json:"vhost_https_port"`
	// TCPMuxHTTPConnectPort 指定服务器侦听TCP HTTP CONNECT请求的端口。如果该值为0，则服务器不会在单个端口上多路传输TCP请求。
	// 如果不是，它将侦听HTTP CONNECT请求的此值。默认情况下，此值为0。
	TCPMuxHTTPConnectPort int `ini:"tcpmux_httpconnect_port" json:"tcpmux_httpconnect_port"`
	// 如果 TCPMuxPassThrough 为true，则frps不会对流量进行任何更新。
	TCPMuxPassThrough bool `ini:"tcpmux_passthrough" json:"tcpmux_passthrough"`
	// VhostHTTPTimeout 指定Vhost HTTP服务器的响应标头超时（以秒为单位）。默认情况下，此值为60。
	VhostHTTPTimeout int64 `ini:"vhost_http_timeout" json:"vhost_http_timeout"`
	// DashboardAddr 指定仪表板绑定到的地址。默认情况下，此值为“0.0.0.0”。
	DashboardAddr string `ini:"dashboard_addr" json:"dashboard_addr"`
	// DashboardPort 指定仪表板侦听的端口。如果此值为 0，则不会启动仪表板。默认情况下，此值为 0。
	DashboardPort int `ini:"Dashboard_port" json:"dashboard_port"`
	// DashboardTLSCertFile 指定服务器将加载的证书文件的路径。
	// 如果“dashboard_tls_cert_file”、“dashboard_tls_key_file”有效，则服务器将使用此提供的 tls 配置。
	DashboardTLSCertFile string `ini:"dashboard_tls_cert_file" json:"dashboard_tls_cert_file"`
	// DashboardTLSKeyFile 指定服务器将加载的密钥的路径。
	// 如果“dashboard_tls_cert_file”、“dashboard_tls_key_file”有效，则服务器将使用此提供的 tls 配置。
	DashboardTLSKeyFile string `ini:"dashboard_tls_key_file" json:"dashboard_tls_key_file"`
	// DashboardTLSMode 指定 HTTP 或 HTTPS 模式之间的仪表板模式。默认情况下，此值为 false，即 HTTP 模式。
	DashboardTLSMode bool `ini:"dashboard_tls_mode" json:"dashboard_tls_mode"`
	// DashboardUser 指定仪表板将用于登录的用户名。
	DashboardUser string `ini:"dashboard_user" json:"dashboard_user"`
	// DashboardPwd 指定仪表板将用于登录的密码。
	DashboardPwd string `ini:"dashboard_pwd" json:"dashboard_pwd"`
	// EnablePrometheus 将在 metrics api 中导出 {dashboard_addr}：{dashboard_port} 上的 prometheus 指标。
	EnablePrometheus bool `ini:"enable_prometheus" json:"enable_prometheus"`
	// AssetsDir 指定仪表板将从中加载资源的本地目录。
	// 如果此值为 ""，则将使用 static 从捆绑的可执行文件中加载资产。默认情况下，此值为“”。
	AssetsDir string `ini:"assets_dir" json:"assets_dir"`
	// LogFile 指定将日志写入的文件。仅当正确设置了 LogWay 时，才会使用此值。
	// 默认情况下，此值为“console”。
	LogFile string `ini:"log_file" json:"log_file"`
	// LogWay 指定日志记录的管理方式。有效值为“console”或“file”。
	// 如果使用“console”，日志将被打印到 stdout。如果使用“file”，日志将打印到 LogFile。
	// 默认情况下，此值为“console”。
	LogWay string `ini:"log_way" json:"log_way"`
	// LogLevel 指定最低日志级别。有效值为 “trace”、“debug”、“info”、“warn” 和 “error”。
	// 默认情况下，此值为“info”。
	LogLevel string `ini:"log_level" json:"log_level"`
	// LogMaxDays 指定删除前存储日志信息的最大天数。这仅在 LogWay == “file” 时使用。
	// 默认情况下，此值为 0。
	LogMaxDays int64 `ini:"log_max_days" json:"log_max_days"`
	// DisableLogColor 当 LogWay == “console” 设置为 true 时禁用日志颜色。
	// 默认情况下，此值为 false。
	DisableLogColor bool `ini:"disable_log_color" json:"disable_log_color"`
	// DetailedErrorsToClient 定义是否将特定错误（包含调试信息）发送到 frpc。
	// 默认情况下，此值为 true。
	DetailedErrorsToClient bool `ini:"detailed_errors_to_client" json:"detailed_errors_to_client"`

	// SubDomainHost 指定在使用 Vhost 代理时将附加到客户端请求的子域的域。
	// 例如，如果此值设置为“frps.com”，并且客户端请求子域“test”，则生成的 URL 将为“test.frps.com”。
	// 默认情况下，此值为“”。
	SubDomainHost string `ini:"sub_domain_host" json:"sub_domain_host"`
	// TCPMux 切换 TCP 流多路复用。这允许来自客户端的多个请求共享单个 TCP 连接。
	// 默认情况下，此值为 true。
	TCPMux bool `ini:"tcp_mux" json:"tcp_mux"`
	// TCPMuxKeepaliveInterval 指定 TCP 流多路复用的保持活动间隔。
	//

	TCPMuxKeepaliveInterval int64 `ini:"tcp_mux_keepalive_interval" json:"tcp_mux_keepalive_interval"`
}

func UnmarshalServerConfFromIni(source interface{})
