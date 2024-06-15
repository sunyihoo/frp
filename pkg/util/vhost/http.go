package vhost

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/fatedier/golib/pool"
	httppkg "github.com/sunyihoo/frp/pkg/util/http"
	"github.com/sunyihoo/frp/pkg/util/log"
	stdlog "log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var ErrNoRouteFound = errors.New("no route found")

type HTTPReverseProxyOptions struct {
	ResponseHeaderTimeoutS int64
}

type HTTPReverseProxy struct {
	proxy       *httputil.ReverseProxy
	vhostRouter *Routers

	responseHeaderTimeout time.Duration
}

func NewHTTPReverseProxy(option HTTPReverseProxyOptions, vhostRouter *Routers) *HTTPReverseProxy {
	if option.ResponseHeaderTimeoutS <= 0 {
		option.ResponseHeaderTimeoutS = 60
	}
	rp := &HTTPReverseProxy{
		responseHeaderTimeout: time.Duration(option.ResponseHeaderTimeoutS) * time.Second,
		vhostRouter:           vhostRouter,
	}
	// todo 学习
	proxy := &httputil.ReverseProxy{
		// 按路由策略修改传入请求。
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]
			r.SetXForwarded()
			req := r.Out
			req.URL.Scheme = "http"
			reqRouteInfo := req.Context().Value(RouteInfoKey).(*RequestRouteInfo)
			originalHost, _ := httppkg.CanonicalHost(reqRouteInfo.Host)

			rc := req.Context().Value(RouteConfigKey).(*RouteConfig)
			if rc != nil {
				if rc.RewriteHost != "" {
					req.Host = rc.RewriteHost
				}

				var endpoint string
				if rc.ChooseEndpointFn != nil {
					// ignore error here, it will CreateConnFn instead later
					endpoint, _ = rc.ChooseEndpointFn()
					reqRouteInfo.Endpoint = endpoint
					log.Tracef("choose endpoint name [%s] for http request host [%s] path [%s] httpuser [%s]",
						endpoint, originalHost, reqRouteInfo.URL, reqRouteInfo.HTTPUser)
				}

				// Set {domain}.{location}.{routeByHttpUser}.{endpoint} as URL host here to let http transport request connections.
				req.URL.Host = rc.Domain + "." +
					base64.StdEncoding.EncodeToString([]byte(rc.Location)) + "." +
					base64.StdEncoding.EncodeToString([]byte(rc.RouteByHTTPUser)) + "." +
					base64.StdEncoding.EncodeToString([]byte(endpoint))

				for k, v := range rc.Headers {
					req.Header.Set(k, v)
				}
			} else {
				req.URL.Host = req.Host
			}
		},
		ModifyResponse: func(r *http.Response) error {
			rc := r.Request.Context().Value(RouteInfoKey).(*RouteConfig)
			if rc != nil {
				for k, v := range rc.ResponseHeaders {
					r.Header.Set(k, v)
				}
			}
			return nil
		},
		// Create a connection to one proxy routed by route policy
		Transport: &http.Transport{
			ResponseHeaderTimeout: rp.responseHeaderTimeout,
			IdleConnTimeout:       60 * time.Second,
			MaxIdleConnsPerHost:   5,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return rp.CreateConnection(ctx.Value(RouteInfoKey).(*RequestRouteInfo), true)
			},
			Proxy: func(req *http.Request) (*url.URL, error) {
				// Use proxy mode if there is host in HTTP first line. 如果 HTTP 第一行中有主机，使用代理模式。
				// GET http://example.com/ HTTP/1.1
				// Host: example.com
				//
				// Normal:
				// GET / HTTP/1.1
				// Host: example.com
				urlHost := req.Context().Value(RouteInfoKey).(*RequestRouteInfo).URLHost
				if urlHost != "" {
					return req.URL, nil
				}
				return nil, nil
			},
		},
		// todo 学习
		BufferPool: pool.NewBuffer(32 * 1024),
		ErrorLog:   stdlog.New(log.NewWriterLogger(log.WarnLevel, 2), "", 0),
		ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Logf(log.WarnLevel, 1, "do http proxy requset [host: %s] error: %v", req.Host, err)
			if err != nil {
				var e net.Error
				if errors.As(err, &e) && e.Timeout() {
					rw.WriteHeader(http.StatusGatewayTimeout)
					return
				}
			}
			rw.WriteHeader(http.StatusNotFound)
			_, _ = rw.Write(getNotFoundPageContent())
		},
	}
	rp.proxy = proxy
	return rp
}

// CreateConnection create a new connection by route config
func (rp *HTTPReverseProxy) CreateConnection(reqRouteInfo *RequestRouteInfo, byEndpoint bool) (net.Conn, error) {
	host, _ := httppkg.CanonicalHost(reqRouteInfo.Host)
	vr, ok := rp.getVhost(host, reqRouteInfo.URL, reqRouteInfo.HTTPUser)
	if ok {
		if byEndpoint {
			fn := vr.payload.(*RouteConfig).CreateConnByEndpointFn
			if fn != nil {
				return fn(reqRouteInfo.Endpoint, reqRouteInfo.RemoteAddr)
			}
		}
		fn := vr.payload.(*RouteConfig).CreateConnFn
		if fn != nil {
			return fn(reqRouteInfo.RemoteAddr)
		}
	}
	return nil, fmt.Errorf("%v: %s %s %s", ErrNoRouteFound, host, reqRouteInfo.URL, reqRouteInfo.HTTPUser)
}

// getVhost tries to get vhost router by route policy
func (rp *HTTPReverseProxy) getVhost(domain, location, routeByHTTPUser string) (*Router, bool) {
	findRouter := func(inDomain, inLocation, inRouteByHTTPUser string) (*Router, bool) {
		vr, ok := rp.vhostRouter.Get(inDomain, inLocation, inRouteByHTTPUser)
		if ok {
			return vr, ok
		}
		// Try to check if there is one proxy that doesn't specify routerByHTTPUser, it means match all.
		vr, ok = rp.vhostRouter.Get(inDomain, inLocation, "")
		if ok {
			return vr, ok
		}
		return nil, false
	}

	// First we check the full hostname
	// if not exist, then check the wildcard_domain such as *.example.com
	vr, ok := findRouter(domain, location, routeByHTTPUser)
	if ok {
		return vr, ok
	}

	// e.g. domain = text.example.com, try to match wildcard domains
	// *.example.com
	// *.com
	domainSplit := strings.Split(domain, ".")
	for {
		if len(domainSplit) < 3 {
			break
		}

		domainSplit[0] = "*"
		domain = strings.Join(domainSplit, ".")
		vr, ok = findRouter(domain, location, routeByHTTPUser)
		if ok {
			return vr, true
		}
		domainSplit = domainSplit[1:]
	}

	// Finally, try to check if there is one proxy that domain is "*" means match all domains.
	vr, ok = findRouter("*", location, routeByHTTPUser)
	if ok {
		return vr, true
	}
	return nil, false
}
