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

package xlog

type LogPrefix struct {
	// Name 是前缀的名称，它不会显示在日志中，而是用于标识前缀。
	Name string
	// Value 是前缀的值，它将显示在日志中。
	Value string
	// 优先级较高的前缀将首先显示，默认值为 10。
	Priority int
}

// Logger 不是前缀操作的线程安全
type Logger struct {
	prefixes []LogPrefix

	prefixString string
}
