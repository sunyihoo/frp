package proxy

import (
	"context"
	"sync"
)

type Proxy interface {
	Context() context.Context
	Run() (remoteAddr string,err error)
	GetName() string
	GetConfigurer() v1.
}

type Manager struct {
	// 按代理名称索引的代理
	pxys map[string]Proxy

	mu sync.Mutex
}
