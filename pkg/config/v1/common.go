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
	"sync"

	"github.com/sunyihoo/frp/pkg/util/util"
)

// TODO(fatedier): Due to the current implementation issue of the go json library, the UnmarshalJSON method
// of a custom struct cannot access the DisallowUnknownFields parameter of the parent decoder.
// Here, a global variable is temporarily used to control whether unknown fields are allowed.
// Once the v2 version is implemented by the community, we can switch to a standardized approach.
//
// https://github.com/golang/go/issues/41144
// https://github.com/golang/go/discussions/63397
var (
	DisallowUnknownFields   = false
	DisallowUnKnownFieldsMu sync.Mutex
)

type AuthScope string

const (
	AuthScopeHeartBeats   AuthScope = "HeartBeats"
	AuthScopeNewWorkConns AuthScope = "NewWorkConns"
)

type AuthMethod string

const (
	AuthMethodToken AuthMethod = "token"
	AuthMethodOIDC  AuthMethod = "oidc"
)

// QUICOptions protocol options
type QUICOptions struct {
	KeepalivePeriod    int `json:"keepalivePeriod,omitempty"`
	MaxIdleTimeout     int `json:"maxIdleTimeout,omitempty"`
	MaxIncomingStreams int `json:"maxIncomingStreams,omitempty"`
}

func (c *QUICOptions) Complete() {
	c.KeepalivePeriod = util.EmptyOr(c.KeepalivePeriod, 10)
	c.MaxIdleTimeout = util.EmptyOr(c.MaxIdleTimeout, 30)
	c.MaxIncomingStreams = util.EmptyOr(c.MaxIncomingStreams, 100000)
}

type WebServerConfig struct {
	// Addr 这是为提供 Web 界面和 API 而绑定的网络地址。
	// 默认情况下，此值为“127.0.0.1”。
	Addr string `json:"addr,omitempty"`
	// Port 端口指定 Web 服务器要侦听的端口。
	// 如果此值为 0，则不会启动管理服务器。
	Port int `json:"port,omitempty"`
	// User 指定 Web 服务器将用于登录的用户名。
	User string `json:"user,omitempty"`
	// Password 指定管理服务器将用于登录的密码。
	Password string `json:"password,omitempty"`
	// AssetsDir指定管理服务器将从中加载资源的本地目录。
	// 如果此值为 “”，则将使用 embed 包从捆绑的可执行文件加载资产。
	AssetsDir string `json:"assetsDir,omitempty"`
	// 启用 golang pprof 处理程序。
	PprofEnable bool `json:"pprofEnable,omitempty"`
	// 如果 TLSConfig 不是 nil，则启用 TLS。
	TLS *TLSConfig `json:"tls,omitempty"`
}

func (c *WebServerConfig) Complete() {
	c.Addr = util.EmptyOr(c.Addr, "127.0.0.1")
}

type TLSConfig struct {
	// CertFile 指定客户端将加载的证书文件的路径。
	CertFile string `json:"certFile,omitempty"`
	// KeyFile 指定客户端将加载的密钥文件的路径。
	KeyFile string `json:"keyFile,omitempty"`
	// TrustedCaFile 指定将加载的受信任 CA 文件的路径。
	TrustedCaFile string `json:"trustedCaFile,omitempty"`
	// ServerName 指定 TLS 证书的自定义服务器名称。
	// 默认情况下，服务器名称（如果与 ServerAddr 相同）。
	ServerName string `json:"serverName,omitempty"`
}

type LogConfig struct {
	// 这是 frp 应该写入日志的目标。
	// 如果使用“控制台”，日志将打印到 stdout，
	// 否则，日志将写入指定文件。
	// 默认情况下，此值为“console”。
	To string `json:"to,omitempty"`
	// Level 指定最低日志级别。有效值为 “trace”、“debug”、“info”、“warn” 和 “error”。
	// 默认情况下，此值为“info”。
	Level string `json:"level,omitempty"`

	// MaxDays 指定删除前存储日志信息的最大天数。
	MaxDays int64 `json:"maxDays,omitempty"`
	// DisabledPrintColor 在“控制台”时禁用日志颜色 log.to。
	DisabledPrintColor bool `json:"disabledPrintColor,omitempty"`
}

func (c *LogConfig) Complete() {
	c.To = util.EmptyOr(c.To, "console")
	c.Level = util.EmptyOr(c.Level, "info")
	c.MaxDays = util.EmptyOr(c.MaxDays, 3)
}

type HTTPPluginOptions struct {
	Name      string   `json:"name"`
	Addr      string   `json:"addr"`
	Path      string   `json:"path"`
	Ops       []string `json:"ops"`
	TLSVerify bool     `json:"tlsVerify,omitempty"`
}

type HeaderOperations struct {
	Set map[string]string `json:"set,omitempty"`
}

type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
