package http

import (
	"crypto/tls"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

type Server struct {
	addr   string
	ln     net.Listener
	tlsCfg *tls.Config

	router *mux.Router
	hs     *http.Server

	authMiddleware mux.MiddlewareFunc
}
