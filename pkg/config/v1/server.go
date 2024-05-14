package v1

type ServerConfig struct {
	APIMetadata

	Auth
}

type AuthServerConfig struct {
	Method           AuthMethod  `json:"method,omitempty"`
	AdditionalScopes []AuthScope `json:"additionalScopes,omitempty"`
	Token            string      `json:"token,omitempty"`
	OIDC
}

type AuthOIDCServerConfig struct {
	// Issuer 指定用于验证OIDC令牌的颁布者
	Issuer string `json:"issuer,omitempty"`
}
