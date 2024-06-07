package vhost

import "sync"

type routerByHTTPUser map[string][]*Router

type Routers struct {
	indexByDomain map[string]routerByHTTPUser

	mutex sync.RWMutex
}

type Router struct {
	domain   string
	location string
	httpUser string

	// 在此处存储任何对象
	payload interface{}
}

func NewRouters() *Routers {
	return &Routers{
		indexByDomain: make(map[string]routerByHTTPUser),
	}

}
