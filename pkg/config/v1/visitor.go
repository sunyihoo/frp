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

type VisitorTransport struct {
	UseEncryption  bool `json:"useEncryption"`
	UseCompression bool `json:"useCompression"`
}

type VisitorBaseConfig struct {
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	Transport VisitorTransport `json:"transport,omitempty"`
	SecretKey string           `json:"secretKey,omitempty"`
	// 如果未设置服务器用户，则默认为当前用户
	ServerUser string `json:"serverUser,omitempty"`
	ServerName string `json:"serverName,omitempty"`
	BindAddr   string `json:"bindAddr,omitempty"`
	// BindPort 是访问者侦听的端口。
	// 它可以小于 0，这意味着不要绑定到端口，只接收从其他访问者重定向的连接。
	// SUDP 现在不支持此功能
	BindPort string `json:"bindPort,omitempty"`
}

type TypeVisitorConfig struct {
	Type string `json:"type"`
	VisitorConfigurer
}

type VisitorConfigurer interface {
	Complete(config *ClientCommonConfig)
	GetBaseConfig() *VisitorBaseConfig
}
