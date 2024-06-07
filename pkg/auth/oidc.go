package auth

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
	"slices"
)

type OidcAuthConsumer struct {
	additionalAuthScopes []v1.AuthScope

	verifier         *oidc.IDTokenVerifier
	subjectFromLogin string
}

func NewOidcAuthVerifier(additionalAuthScopes []v1.AuthScope, cfg v1.AuthOIDCServerConfig) *OidcAuthConsumer {
	provider, err := oidc.NewProvider(context.Background(), cfg.Issuer)
	if err != nil {
		panic(err)
	}
	verifierConf := oidc.Config{
		ClientID:          cfg.Audience,
		SkipClientIDCheck: cfg.Audience == "",
		SkipExpiryCheck:   cfg.SkipExpiryCheck,
		SkipIssuerCheck:   cfg.SkipIssuerCheck,
	}
	return &OidcAuthConsumer{
		additionalAuthScopes: additionalAuthScopes,
		verifier:             provider.Verifier(&verifierConf),
	}
}

func (auth *OidcAuthConsumer) VerifyLogin(loginMsg *msg.Login) (err error) {
	token, err := auth.verifier.Verify(context.Background(), loginMsg.PrivilegeKey)
	if err != nil {
		return fmt.Errorf("invalid OIDC token in login: %v", err)
	}
	auth.subjectFromLogin = token.Subject
	return nil
}

func (auth *OidcAuthConsumer) verifyPostLoginToken(privilegeKey string) (err error) {
	token, err := auth.verifier.Verify(context.Background(), privilegeKey)
	if err != nil {
		return fmt.Errorf("invalid OIDC token in ping: %v", err)
	}
	if token.Subject != auth.subjectFromLogin {
		return fmt.Errorf("recevied different OIDC subject in login and ping. "+
			"original subject: %s, "+
			"new subject: %s",
			auth.subjectFromLogin, token.Subject)
	}
	return nil
}

func (auth *OidcAuthConsumer) VerifyPing(pingMsg *msg.Ping) (err error) {
	if !slices.Contains(auth.additionalAuthScopes, v1.AuthScopeHeartBeats) {
		return nil
	}

	return auth.verifyPostLoginToken(pingMsg.PrivilegeKey)
}

func (auth *OidcAuthConsumer) VerifyNewWorkConn(newWorkConn *msg.NewWorkConn) (err error) {
	if !slices.Contains(auth.additionalAuthScopes, v1.AuthScopeNewWorkConns) {
		return nil
	}

	return auth.verifyPostLoginToken(newWorkConn.PrivilegeKey)
}
