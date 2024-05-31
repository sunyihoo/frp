package legacy

type BaseConfig struct {
	// AuthenticationMethod 指定使用何种身份验证方法对frpc和frps进行身份验证。
	// 如果指定了“token” - 则token将被读取到登录消息中。
	// 如果指定了“oidc” - 将使用oidc设置颁发oidc（Open ID Connect）令牌。
	// 默认情况下，此值为 "token"。
	AuthenticationMethod string `ini:"authentication_method" json:"authentication_method"`
	// AuthenticateHeartBeats 指定是否在发送到frps的检测信号中包括身份验证令牌。
	// 默认情况下，此值为 false
	AuthenticateHeartBeats bool `ini:"authenticate_heartbeats" json:"authenticate_heartbeats"`
	// AuthenticateNewWorkConns 指定是否在发送到frps的新工作连接中包括身份验证令牌。
	// 默认情况下，此值为 false
	AuthenticateNewWorkConns bool `ini:"authenticate_new_work_conns" json:"authenticate_new_work_conns"`
}

func getDefaultBaseConf() BaseConfig {
	return BaseConfig{
		AuthenticationMethod:     "token",
		AuthenticateHeartBeats:   false,
		AuthenticateNewWorkConns: false,
	}
}

type ServerConfig struct {
	BaseConfig       `ini:",extends"`
	OidcServerConfig `ini:",extends"`
	TokenConfig      `ini:",extends"`
}

func GetDefaultServerConf() ServerConfig {
	return ServerConfig{
		BaseConfig:       getDefaultBaseConf(),
		OidcServerConfig: getDefaultOidcServerConf(),
		TokenConfig:      getDefaultTokenConf(),
	}
}

type OidcServerConfig struct {
	// OidcIssuer 指定用于验证OIDC令牌的颁发者。此颁发者将用于加载公钥以验证签名，并将与OIDC令牌中的颁发者声明进行比较。
	// 如果 AuthenticationMethod == "oidc"，将使用它。默认情况下，此值为""。
	OidcIssuer string `ini:"oidc_issuer" json:"oidc_issuer"`
	// OidcAudience 指定验证时OIDC令牌应包含的受众。如果此值为空，则将跳过访问群体（"client ID"）验证。
	// 它将在AuthenticationMethod==“oidc”时使用。默认情况下，此值为""。
	OidcAudience string `ini:"oidc_audience" json:"oidc_audience"`
	// OidcSkipExpiryCheck 指定如果OIDC令牌已过期，是否跳过检查。
	// 它将在 AuthenticationMethod == "oidc" 时使用。默认情况下，此值为false。
	OidcSkipExpiryCheck bool `ini:"oidc_skip_expiry_check" json:"oidc_skip_expiry_check"`
	// OidcSkipIssuerCheck 指定是否跳过检查OIDC令牌的颁发者声明是否与 OidcIssuer 中指定的颁发者匹配。
	// 它将在 AuthenticationMethod == "oidc"时使用。默认情况下，此值为false。
	OidcSkipIssuerCheck bool `ini:"oidc_skip_issuer_check" json:"oidc_skip_issuer_check"`
}

func getDefaultOidcServerConf() OidcServerConfig {
	return OidcServerConfig{
		OidcIssuer:          "",
		OidcAudience:        "",
		OidcSkipIssuerCheck: false,
		OidcSkipExpiryCheck: false,
	}
}

type TokenConfig struct {
	// Token 令牌指定用于创建要发送到服务器的密钥的授权令牌。
	// 服务器必须具有匹配的令牌，授权才能成功。默认情况下，此值为 ""。
	Token string `ini:"token" json:"token"`
}

func getDefaultTokenConf() TokenConfig {
	return TokenConfig{
		Token: "",
	}
}
