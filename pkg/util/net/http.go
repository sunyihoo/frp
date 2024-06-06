package net

import (
	"github.com/sunyihoo/frp/pkg/util/util"
	"net/http"
	"time"
)

type HTTPAuthMiddleware struct {
	user          string
	passwd        string
	authFailDelay time.Duration
}

func NewHTTPAuthMiddleware(user, passwd string) *HTTPAuthMiddleware {
	return &HTTPAuthMiddleware{
		user:   user,
		passwd: passwd,
	}
}

func (authMid *HTTPAuthMiddleware) SetAuthFailDelay(delay time.Duration) *HTTPAuthMiddleware {
	authMid.authFailDelay = delay
	return authMid
}

func (authMid *HTTPAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqUser, reqPasswd, hasAuth := r.BasicAuth()
		if (authMid.user == "" && authMid.passwd == "") ||
			(hasAuth && util.ConstantTimeEqString(reqUser, authMid.user) &&
				util.ConstantTimeEqString(reqPasswd, authMid.passwd)) {
			next.ServeHTTP(w, r)
		} else {
			if authMid.authFailDelay > 0 {
				time.Sleep(authMid.authFailDelay)
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}
