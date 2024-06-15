// Copyright 2019 fatedier, fatedier@gmail.com
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

package server

type Manager struct {
	loginPlugins       []Plugin
	newProxyPlugins    []Plugin
	closeProxyPlugins  []Plugin
	pingPlugins        []Plugin
	newWorkConnPlugins []Plugin
	newUserConnPlugins []Plugin
}

func NewManager() *Manager {
	return &Manager{
		loginPlugins:       make([]Plugin, 0),
		newProxyPlugins:    make([]Plugin, 0),
		closeProxyPlugins:  make([]Plugin, 0),
		pingPlugins:        make([]Plugin, 0),
		newWorkConnPlugins: make([]Plugin, 0),
		newUserConnPlugins: make([]Plugin, 0),
	}
}

func (m *Manager) Register(p Plugin) {
	if p.IsSupport(OpLogin) {
		m.loginPlugins = append(m.loginPlugins, p)
	}
	if p.IsSupport(OpNewProxy) {
		m.newProxyPlugins = append(m.newProxyPlugins, p)
	}
	if p.IsSupport(OpCloseProxy) {
		m.closeProxyPlugins = append(m.closeProxyPlugins, p)
	}
	if p.IsSupport(OpPing) {
		m.pingPlugins = append(m.pingPlugins, p)
	}
	if p.IsSupport(OpNewWorkConn) {
		m.newWorkConnPlugins = append(m.newWorkConnPlugins, p)
	}
	if p.IsSupport(OpNewUserConn) {
		m.newUserConnPlugins = append(m.newUserConnPlugins, p)
	}
}
