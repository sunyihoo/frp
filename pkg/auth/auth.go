package auth

import (
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
)

type Setter interface {
	SetLogin(*msg.Login) error
	SetPing(*msg.Ping) error
	SetNewWorkConn(*msg.NewWorkConn) error
}

type Verifier interface {
	VerifyLogin(*msg.Login) error
	VerifyPing(*msg.Ping) error
	VerifyNewWorkConn(conn *msg.NewWorkConn) error
}

func NewAuthVerifier(cfg v1.AuthServerConfig) (authVerifier Verifier) {
	switch cfg.Method {
	case v1.AuthMethodToken:
		authVerifier = NewTokenAuth(cfg.AdditionalScopes, cfg.Token)
	case v1.AuthMethodOIDC:
		authVerifier = NewOidcAuthVerifier(cfg.AdditionalScopes, cfg.OIDC)
	}
	return authVerifier
}
