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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/samber/lo"

	"github.com/sunyihoo/frp/pkg/config/types"
	"github.com/sunyihoo/frp/pkg/msg"
	"github.com/sunyihoo/frp/pkg/util/util"
)

type ProxyTransport struct {
	// UseEncryption 控制是否对与服务器的通信进行加密。
	// 加密是使用服务器和客户端配置中提供的令牌完成的。
	//
	UseEncryption bool `json:"useEncryption,omitempty"`

	// UseCompression 控制是否压缩与服务器的通信。
	UseCompression bool `json:"useCompression,omitempty"`
	// BandwidthLimit 限制带宽
	// 0 意味着不限制
	BandwidthLimit types.BandwidthQuantity `json:"bandwidthLimit,omitempty"`
	// BandwidthLimitMode 指定是限制客户端还是服务器端的带宽。
	// 有效值包括“client”和“server”。
	// 默认情况下，此值为“client”。
	BandwidthLimitMode string `json:"bandwidthLimitMode,omitempty"`
	// ProxyProtocolVersion 指定要使用的协议版本。
	// 有效值包括"v1"、"v2"和""。
	// 如果值为""，则将自动选择协议版本。默认情况下，此值为""。
	ProxyProtocolVersion string `json:"proxyProtocolVersion,omitempty"`
}

type LoadBalanceConfig struct {
	// Group 指定所属的组。服务器将使用此信息
	// 对同一组中的代理进行负载平衡。
	// 如果值为 ""，则该值将不在组中。
	Group string `json:"group"`
	// GroupKey 指定一个组密钥，
	// 该密钥在同一组的代理之间应相同。
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

// HealthCheckConfig 配置运行状况检查。这对于负载平衡目的非常有用，
// 以检测和删除故障服务的代理。
type HealthCheckConfig struct {
	// Type 指定用于运行状况检查的协议。
	// 有效值包括"tcp"、"http" 和 ""。
	// 如果该值为""，则不会进行健康检查。
	//
	// 如果类型为“tcp”，则将尝试与目标服务器建立连接。
	// 如果无法建立连接，则健康检查失败。
	//
	// 如果类型为“http”，则将对 HealthCheckURL 指定的
	// 终结点发出 GET 请求。
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
	// 如果运行状况检查类型为“http”，
	// 则 Path 指定要将运行状况检查发送到的路径。
	Path string `json:"path,omitempty"`
	// 如果运行状况检查类型为“http”，
	// 则 HTTPHeaders 指定要与运行状况请求一起发送的请求头。
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

func (c *ProxyBaseConfig) GetBaseConfig() *ProxyBaseConfig {
	return c
}

func (c *ProxyBaseConfig) Complete(namePrefix string) {
	// todo 学习
	c.Name = lo.Ternary(namePrefix == "", "", namePrefix+".") + c.Name
	c.LocalIP = util.EmptyOr(c.LocalIP, "127.0.0.1")
	c.Transport.BandwidthLimitMode = util.EmptyOr(c.Transport.BandwidthLimitMode, types.BandwidthLimitModeClient)
}

func (c *ProxyBaseConfig) MarshalToMsg(m *msg.NewProxy) {
	m.ProxyName = c.Name
	m.ProxyType = c.Type
	m.UseEncryption = c.Transport.UseEncryption
	m.UseCompression = c.Transport.UseCompression
	m.BandwidthLimit = c.Transport.BandwidthLimit.String()
	// 默认值留空以减少流量
	if c.Transport.BandwidthLimitMode != "client" {
		m.BandwidthLimitMode = c.Transport.BandwidthLimitMode
	}
	m.Group = c.LoadBalancer.Group
	m.GroupKey = c.LoadBalancer.GroupKey
	m.Metas = c.Metadatas
	m.Annotations = c.Annotations
}

func (c *ProxyBaseConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.Name = m.ProxyName
	c.Type = m.ProxyType
	c.Transport.UseEncryption = m.UseEncryption
	c.Transport.UseCompression = m.UseCompression
	if m.BandwidthLimit != "" {
		c.Transport.BandwidthLimit, _ = types.NewBandwidthQuantity(m.BandwidthLimit)
	}
	if m.BandwidthLimitMode != "" {
		c.Transport.BandwidthLimitMode = m.BandwidthLimitMode
	}
	c.LoadBalancer.Group = m.Group
	c.LoadBalancer.GroupKey = m.GroupKey
	c.Metadatas = m.Metas
	c.Annotations = m.Annotations
}

type TypedProxyConfig struct {
	Type string `json:"type"`
	ProxyConfigurer
}

func (c *TypedProxyConfig) UnmarshalJSON(b []byte) error {
	if len(b) == 4 && string(b) == "null" {
		return errors.New("type is required")
	}

	typeStruct := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(b, &typeStruct); err != nil {
		return err
	}

	c.Type = typeStruct.Type
	configurer := NewProxyConfigurerByType(ProxyType(c.Type))
	if configurer == nil {
		return fmt.Errorf("unknown proxy type: %s", c.Type)
	}
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	if DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(configurer); err != nil {
		return fmt.Errorf("unmarshal ProxyConfig error: %s", err)
	}
	c.ProxyConfigurer = configurer
	return nil
}

type ProxyConfigurer interface {
	Complete(namePrefix string)
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
	ProxyTypeTCP    ProxyType = "tcp"
	ProxyTypeUDP    ProxyType = "udp"
	ProxyTypeTCPMUX ProxyType = "tcpmux"
	ProxyTypeHTTP   ProxyType = "http"
	ProxyTypeHTTPS  ProxyType = "https"
	ProxyTypeSTCP   ProxyType = "stcp"
	ProxyTypeXTCP   ProxyType = "xtcp"
	ProxyTypeSUDP   ProxyType = "sudp"
)

var proxyConfigTypeMap = map[ProxyType]reflect.Type{
	ProxyTypeTCP:    reflect.TypeOf(TCPProxyConfig{}),
	ProxyTypeUDP:    reflect.TypeOf(UDPProxyConfig{}),
	ProxyTypeHTTP:   reflect.TypeOf(HTTPProxyConfig{}),
	ProxyTypeHTTPS:  reflect.TypeOf(HTTPSProxyConfig{}),
	ProxyTypeTCPMUX: reflect.TypeOf(TCPMuxProxyConfig{}),
	ProxyTypeSTCP:   reflect.TypeOf(STCPProxyConfig{}),
	ProxyTypeXTCP:   reflect.TypeOf(XTCPProxyConfig{}),
	ProxyTypeSUDP:   reflect.TypeOf(UDPProxyConfig{}),
}

func NewProxyConfigurerByType(proxyType ProxyType) ProxyConfigurer {
	v, ok := proxyConfigTypeMap[proxyType]
	if !ok {
		return nil
	}
	pc := reflect.New(v).Interface().(ProxyConfigurer)
	pc.GetBaseConfig().Type = string(proxyType)
	return pc
}

var _ ProxyConfigurer = &TCPProxyConfig{}

type TCPProxyConfig struct {
	ProxyBaseConfig

	RemotePort int `json:"remotePort,omitempty"`
}

func (c *TCPProxyConfig) MarshalToMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.RemotePort = c.RemotePort
}

func (c *TCPProxyConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	c.RemotePort = m.RemotePort
}

var _ ProxyConfigurer = &UDPProxyConfig{}

type UDPProxyConfig struct {
	ProxyBaseConfig

	RemotePort int `json:"remotePort,omitempty"`
}

func (c *UDPProxyConfig) MarshalToMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.RemotePort = c.RemotePort
}

func (c *UDPProxyConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	c.RemotePort = m.RemotePort
}

var _ ProxyConfigurer = &HTTPProxyConfig{}

type HTTPProxyConfig struct {
	ProxyBaseConfig
	DomainConfig

	Locations         []string         `json:"locations,omitempty"`
	HTTPUser          string           `json:"httpUser,omitempty"`
	HTTPPassword      string           `json:"httpPassword,omitempty"`
	HostHeaderRewrite string           `json:"hostHeaderRewrite,omitempty"`
	RequestHeaders    HeaderOperations `json:"requestHeaders,omitempty"`
	ResponseHeaders   HeaderOperations `json:"responseHeaders,omitempty"`
	RouteByHTTPUser   string           `json:"routeByHTTPUser,omitempty"`
}

func (c *HTTPProxyConfig) MarshalToMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.CustomDomains = c.CustomDomains
	m.SubDomain = c.SubDomain
	m.Locations = c.Locations
	m.HostHeaderRewrite = c.HostHeaderRewrite
	m.HTTPUser = c.HTTPUser
	m.HTTPPwd = c.HTTPPassword
	m.Headers = c.RequestHeaders.Set
	m.ResponseHeaders = c.ResponseHeaders.Set
	m.RouteByHTTPUser = c.RouteByHTTPUser
}

func (c *HTTPProxyConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	c.CustomDomains = m.CustomDomains
	c.SubDomain = m.SubDomain
	c.Locations = m.Locations
	c.HostHeaderRewrite = m.HostHeaderRewrite
	c.HTTPUser = m.HTTPUser
	c.HTTPPassword = m.HTTPPwd
	c.RequestHeaders.Set = m.Headers
	c.ResponseHeaders.Set = m.ResponseHeaders
	c.RouteByHTTPUser = m.RouteByHTTPUser
}

var _ ProxyConfigurer = &HTTPSProxyConfig{}

type HTTPSProxyConfig struct {
	ProxyBaseConfig
	DomainConfig
}

func (c *HTTPSProxyConfig) MarshalToMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.CustomDomains = c.CustomDomains
	m.SubDomain = c.SubDomain
}

func (c *HTTPSProxyConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	c.CustomDomains = m.CustomDomains
	c.SubDomain = m.SubDomain
}

type TCPMultiplexerType string

const (
	TCPMultiplexerHTTPConnect TCPMultiplexerType = "httpconnect"
)

var _ ProxyConfigurer = &TCPMuxProxyConfig{}

type TCPMuxProxyConfig struct {
	ProxyBaseConfig
	DomainConfig

	HTTPUser        string `json:"httpUser,omitempty"`
	HTTPPassword    string `json:"httpPassword,omitempty"`
	RouteByHTTPUser string `json:"routeByHTTPUser,omitempty"`
	Multiplexer     string `json:"multiplexer,omitempty"`
}

func (c *TCPMuxProxyConfig) MarshalToMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.CustomDomains = c.CustomDomains
	m.SubDomain = c.SubDomain
	m.HTTPUser = c.HTTPUser
	m.HTTPPwd = c.HTTPPassword
	m.RouteByHTTPUser = c.RouteByHTTPUser
	m.Multiplexer = c.Multiplexer
}

func (c *TCPMuxProxyConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	c.CustomDomains = m.CustomDomains
	c.SubDomain = m.SubDomain
	c.HTTPUser = m.HTTPUser
	c.HTTPPassword = m.HTTPPwd
	c.RouteByHTTPUser = m.RouteByHTTPUser
	c.Multiplexer = m.Multiplexer
}

var _ ProxyConfigurer = &STCPProxyConfig{}

type STCPProxyConfig struct {
	ProxyBaseConfig

	Secretkey  string   `json:"secretKey,omitempty"`
	AllowUsers []string `json:"allowUsers,omitempty"`
}

func (c *STCPProxyConfig) MarshalToMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.Sk = c.Secretkey
	m.AllowUsers = c.AllowUsers
}

func (c *STCPProxyConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	c.Secretkey = m.Sk
	c.AllowUsers = m.AllowUsers
}

var _ ProxyConfigurer = &XTCPProxyConfig{}

type XTCPProxyConfig struct {
	ProxyBaseConfig

	Secretkey  string   `json:"secretKey,omitempty"`
	AllowUsers []string `json:"allowUsers,omitempty"`
}

func (c *XTCPProxyConfig) MarshalToMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.MarshalToMsg(m)

	m.Sk = c.Secretkey
	m.AllowUsers = c.AllowUsers
}

func (c *XTCPProxyConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	c.Secretkey = m.Sk
	c.AllowUsers = m.AllowUsers
}

var _ ProxyConfigurer = &SUDPProxyConfig{}

type SUDPProxyConfig struct {
	ProxyBaseConfig

	Secretkey  string   `json:"secretKey,omitempty"`
	AllowUsers []string `json:"allowUsers,omitempty"`
}

func (c *SUDPProxyConfig) MarshalToMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	m.Sk = c.Secretkey
	m.AllowUsers = c.AllowUsers
}

func (c *SUDPProxyConfig) UnmarshalFromMsg(m *msg.NewProxy) {
	c.ProxyBaseConfig.UnmarshalFromMsg(m)

	c.Secretkey = m.Sk
	m.AllowUsers = c.AllowUsers
}
