package vhost

import (
	"net/http/httputil"
	"time"
)

type HTTPReverseProxy struct {
	proxy       *httputil.ReverseProxy
	vhostRouter *Routers

	responseHeaderTimeout time.Duration
}
