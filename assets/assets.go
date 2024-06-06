// Copyright 2016 fatedier, fatedier@gmail.com
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

package assets

import (
	"io/fs"
	"net/http"
)

var (
	// 通过“embed”为嵌入式文件创建的只读文件系统
	content fs.FS

	FileSystem http.FileSystem

	// 如果前缀不为空，则从磁盘获取文件内容 content
	prefixPath string
)

// Load 如果路径为空，则在内存中加载资产
// 或使用磁盘文件设置 FileSystem 。
func Load(path string) {
	prefixPath = path
	if prefixPath != "" {
		FileSystem = http.Dir(prefixPath)
	} else {
		FileSystem = http.FS(content)
	}
}
