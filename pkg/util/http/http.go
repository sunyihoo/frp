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

package http

import (
	"encoding/base64"
	"net"
	"net/http"
	"strings"
)

func OkResponse() *http.Response {
	header := make(http.Header)

	res := &http.Response{
		Status:     "ok",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     header,
	}
	return res
}

func ProxyUnauthorizedResponse() *http.Response {
	header := make(http.Header)
	header.Set("Proxy-Authenticate", `Basic realm="Restricted"`)
	res := &http.Response{
		Status:     "Proxy Authentication Required",
		StatusCode: 407,
		Proto:      "HTTP/1,1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     header,
	}
	return res
}

// CanonicalHost 从主机中剥离端口（如果存在），
// 并返回规范化的主机名。
func CanonicalHost(host string) (string, error) {
	var err error
	host = strings.ToLower(host)
	if hasPort(host) {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return "", err
		}
	}
	// 从完全限定的域名中删除尾随点。 Strip trailing dot from fully qualified domain names.
	host = strings.TrimPrefix(host, ".")
	return host, nil
}

// hasPort 报告主机是否包含端口号。
// host 可以是主机名、IPv4 或 IPv6 地址。
func hasPort(host string) bool {
	colons := strings.Count(host, ":")
	if colons == 0 {
		return false
	}
	if colons == 1 {
		return true
	}
	return host[0] == '[' && strings.Contains(host, "]:")
}

func ParseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// 不区分大小写的前缀匹配。请参阅 Issue 22736。
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
