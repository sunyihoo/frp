package legacy

import (
	legacyauth "github.com/sunyihoo/frp/pkg/auth/legacy"
	"gopkg.in/ini.v1"
)

type HTTPPluginOptions struct {
	Name      string `ini:"name"`
	Addr      string `ini:"addr"`
	Path      string `ini:"path"`
	Ops       string `ini:"ops"`
	TLSVerify bool   `ini:"tlsVerify"`
}

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
	// 如果 TCPMux 为true，则不需要应用层的心跳，因为它只能依赖于 TCPMux 中的心跳
	TCPMuxKeepaliveInterval int64 `ini:"tcp_mux_keepalive_interval" json:"tcp_mux_keepalive_interval"`
	// TCPKeepAlive 指定 frpc 和 frps 之间的活动网络连接的保持活动探测器之间的存活时间。
	// 如果为负数，则禁用keep-alive probes。
	TCPKeepAlive int64 `ini:"tcp_keepalive" json:"tcp_keepalive"`
	// Custom404Page 指定要显示的自定义 404 页的路径。
	// 如果此值为 “”，则将显示默认页面。默认情况下，此值为“”。
	Custom404Page string `ini:"custom_404_page" json:"custom_404_page"`

	// AllowPorts 指定客户端能够代理到的一组端口。如果此值的长度为 0，则允许所有端口。
	// 默认情况下，此值为空集。
	AllowPorts map[int]struct{} `ini:"-" json:"-"`
	// Original string
	AllowPortsStr string `ini:"-" json:"-"`
	// MaxPoolCount 指定每个代理的最大池大小。默认情况下，此值为 5。
	MaxPoolCount int64 `ini:"max_pool_count" json:"max_pool_count"`
	// MaxPortsPerClient 指定单个客户端可以代理到的最大端口数。
	// 如果此值为 0，则不会应用任何限制。默认情况下，此值为 0。
	MaxPortsPerClient int64 `ini:"max_ports_per_client" json:"max_ports_per_client"`
	// TLSOnly 指定是否仅接受 TLS 加密的连接。默认情况下，该值为 false。
	TLSOnly bool `ini:"tls_only" json:"tls_only"`
	// TLSCertFile 指定服务器将加载的证书文件的路径。如果“tls_cert_file”、“tls_key_file”有效，
	//则服务器将使用此提供的 tls 配置。否则，服务器将使用自身生成的 tls 配置。
	TLSCertFile string `ini:"tls_cert_file" json:"tls_cert_file"`
	// TLSKeyFile 指定服务器将加载的密钥的路径。如果“tls_cert_file”、“tls_key_file”有效，
	// 则服务器将使用此提供的 tls 配置。否则，服务器将使用自身生成的 tls 配置。
	TLSKeyFile string `ini:"tls_key_file" json:"tls_key_file"`
	// TLSTrustedCaFile 指定服务器将加载的客户端证书文件的路径。它仅在“tls_only”为真时才有效。
	// 如果“tls_trusted_ca_file”有效，服务器将验证每个客户端的证书。
	TLSTrustedCaFile string `ini:"tls_trusted_ca_file" json:"tls_trusted_ca_file"`
	// HeartbeatTimeout 指定在终止连接之前等待检测信号的最长时间。
	// 不建议更改此值。默认情况下，此值为 90。设置负值以禁用它。
	HeartbeatTimeout int64 `ini:"heartbeat_timeout" json:"heartbeat_timeout"`
	// UseConnTimeout 指定等待工作连接的最长时间。默认情况下，此值为 10。
	UserConnTimeout int64 `ini:"user_conn_timeout" json:"user_conn_timeout"`
	// HTTPPlugins
	HTTPPlugins map[string]HTTPPluginOptions `ini:"http_conn_timeout" json:"http_conn_timeout"`
	// UDPPacketSize 指定 UDP 数据包大小 默认情况下，此值为 1500。
	UDPPacketSize int64 `ini:"udp_packet_size" json:"udp_packet_size"`
	// PprofEnabled 在仪表板侦听器中启用 golang pprof 处理程序。
	// 必须先设置仪表板端口。
	PprofEnabled bool `ini:"pprof_enabled" json:"pprof_enabled"`
	// NatHoleAnalysisDataReserveHours 指定 reserve nat hole analysis data的小时数。
	NatHoleAnalysisDataReserveHours int64 `ini:"nat_hole_analysis_data_reserve_hours" json:"nat_hole_analysis_data_reserve_hours"`
}

// GetDefaultServerConf 返回具有合理默认值的服务器配置。
func GetDefaultServerConf() ServerCommonConf {
	return ServerCommonConf{
		ServerConfig: legacyauth.GetDefaultServerConf(),
	}
}

func UnmarshalServerConfFromIni(source interface{}) (ServerCommonConf, error) {
	f, err := ini.LoadSources(ini.LoadOptions{
		Insensitive:         false,
		InsensitiveSections: false,
		InsensitiveKeys:     false,
		IgnoreInlineComment: true,
		AllowBooleanKeys:    true,
	}, source)

	if err != nil {
		return ServerCommonConf{}, err
	}

	s, err := f.GetSection("common")
	if err != nil {
		return ServerCommonConf{}, err
	}

}
