package v1

import (
	"github.com/sunyihoo/frp/pkg/config/types"
	"github.com/sunyihoo/frp/pkg/msg"
)

type ProxyTransport struct {
	// UseEncryption 控制是否对与服务器的通信进行加密。
	// 加密是使用服务器和客户端配置中提供的令牌完成的。
	UseEncryption bool `json:"useEncryption,omitempty"`
	// UseCompression 控制是否压缩与服务器的通信。
	UseCompression bool `json:"useCompression,omitempty"`
	// BandwidthLimit 限制带宽
	// 0 意味着不限制
	BandwidthLimit types.BandwidthWithQuantity `json:"bandwidthLimit,omitempty"`
	// BandwidthLimitMode 指定是限制客户端还是服务器端的带宽。
	// 有效值包括“client”和“server”。默认情况下，此值为“client”。
	BandwidthLimitMode string `json:"bandwidthLimitMode,omitempty"`
	// ProxyProtocolVersion 指定要使用的协议版本。有效值包括"v1"、"v2"和""。
	// 如果值为""，则将自动选择协议版本。默认情况下，此值为""。
	ProxyProtocolVersion string `json:"proxyProtocolVersion,omitempty"`
}

type LoadBalanceConfig struct {
	// Group 指定所属的组。服务器将使用此信息对同一组中的代理进行负载平衡。
	// 如果值为 ""，则该值将不在组中。
	Group string `json:"group"`
	// GroupKey 指定一个组密钥，该密钥在同一组的代理之间应相同。
	GroupKey string `json:"groupKey"`
}

type ProxyBackend struct {
	// LocalIP 指定后端的 IP 地址或主机名。
	LocalIP string `json:"localIP,omitempty"`
	// LocalPort 指定后端的端口。
	LocalPort int `json:"localPort,omitempty"`

	// Plugin 指定应使用哪个插件来处理连接。
	// 如果设置了此值，则将忽略 LocalIP 和 LocalPort 值。
	Plugin TypedClientPluginOptions
}

type HealthCheckConfig struct {
	// Type 指定用于运行状况检查的协议。
	// 有效值包括"tcp"、"http" 和 ""。
	// 如果该值为""，则不会进行健康检查。
	//
	// 如果类型为“tcp”，则将尝试与目标服务器建立连接。
	// 如果无法建立连接，则健康检查失败。
	//
	// 如果类型为“http”，则将对 HealthCheckURL 指定的终结点发出 GET 请求。
	// 如果响应不是 200，则运行状况检查失败。
	Type string `json:"type"` // tcp | http
	// TimeoutSeconds 指定等待运行状况检查尝试连接的秒数。
	// 如果达到超时，则计为运行状况检查失败。
	// 默认情况下，此值为 3。
	TimeoutSeconds int `json:"timeoutSeconds,omitempty"`
	// MaxFailed 指定停止之前允许的失败数。
	// 默认情况下，此值为 1。
	MaxFailed int `json:"maxFailed,omitempty"`
	// IntervalSeconds 指定运行状况检查之间的时间（以秒为单位）。
	// 默认情况下，此值为 10。
	IntervalSeconds int `json:"intervalSeconds"`
	// 如果运行状况检查类型为“http”，则 Path 指定要将运行状况检查发送到的路径。
	Path string `json:"path,omitempty"`
	// 如果运行状况检查类型为“http”，则 HTTPHeaders 指定要与运行状况请求一起发送的请求头。
	HTTPHeaders []HTTPHeader `json:"HTTPHeaders,omitempty"`
}

type DomainConfig struct {
	CustomDomains []string `json:"customDomains,omitempty"`
	SubDomain     string   `json:"subDomain,omitempty"`
}

type ProxyBaseConfig struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Transport   ProxyTransport    `json:"transport,omitempty"`
	// 每个代理的元数据信息
	Metadatas    map[string]string `json:"metadatas,omitempty"`
	LoadBalancer LoadBalanceConfig `json:"loadBalancer,omitempty"`
	HealthCheck  HealthCheckConfig `json:"healthCheck,omitempty"`
	ProxyBackend
}

type ProxyConfigurer interface {
	Complete() (namePrefix string)
	GetBaseConfig() *ProxyBaseConfig
	// MarshalToMsg 将此配置序列化成 msg.NewProxy 消息。
	// 此函数将在 frpc 端调用。
	MarshalToMsg(*msg.NewProxy)
	// UnmarshalFromMsg 反序列化成 msg.NewProxy 消息添加到此配置中。
	// 此函数将在 frps 端调用。
	UnmarshalFromMsg(*msg.NewProxy)
}

type ProxyType string

const (
	ProxyTypeTcp    ProxyType = "tcp"
	ProxyTypeUDP    ProxyType = "udp"
	ProxyTypeTCPMUX ProxyType = "tcpmux"
	ProxyTypeHTTP   ProxyType = "http"
	ProxyTypeHTTPS  ProxyType = "https"
	ProxyTypeSTCP   ProxyType = "stcp"
	ProxyTypeXTCP   ProxyType = "xtcp"
	ProxyTypeSUDP   ProxyType = "sudp"
)
