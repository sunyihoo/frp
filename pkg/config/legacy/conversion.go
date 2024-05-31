package legacy

import (
	"github.com/samber/lo"
	"github.com/sunyihoo/frp/pkg/config/types"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
)

func Convert_ServerCommonConf_To_v1(conf *ServerCommonConf) *v1.ServerConfig {
	out := &v1.ServerConfig{}
	out.Auth.Method = v1.AuthMethod(conf.ServerConfig.AuthenticationMethod)
	out.Auth.Token = conf.ServerConfig.Token
	if conf.ServerConfig.AuthenticateHeartBeats {
		out.Auth.AdditionalScopes = append(out.Auth.AdditionalScopes, v1.AuthScopeHeartBeats)
	}
	if conf.ServerConfig.AuthenticateNewWorkConns {
		out.Auth.AdditionalScopes = append(out.Auth.AdditionalScopes, v1.AuthScopeNewWorkConns)
	}
	out.Auth.OIDC.Audiences = conf.ServerConfig.OidcAudience
	out.Auth.OIDC.Issuer = conf.ServerConfig.OidcIssuer
	out.Auth.OIDC.SkipExpiryCheck = conf.ServerConfig.OidcSkipExpiryCheck
	out.Auth.OIDC.SkipIssuerCheck = conf.ServerConfig.OidcSkipIssuerCheck

	out.BindAddr = conf.BindAddr
	out.BindPort = conf.BindPort
	out.KCPBindPort = conf.KCPBindPort
	out.QUICBindPort = conf.QUICBindPort
	out.Transport.QUIC.KeepalivePeriod = conf.QUICKeepalivePeriod
	out.Transport.QUIC.MaxIdleTimeout = conf.QUICMaxIdleTimeout
	out.Transport.QUIC.MaxIncomingStreams = conf.QUICMaxIncomingStreams

	out.ProxyBindAddr = conf.ProxyBindAddr
	out.VhostHTTPPort = conf.VhostHTTPPort
	out.VhostHTTPSPort = conf.VhostHTTPSPort
	out.TCPMuxHTTPConnectPort = conf.TCPMuxHTTPConnectPort
	out.TCPMuxPassthrough = conf.TCPMuxPassthrough
	out.VhostHTTPTimeout = conf.VhostHTTPTimeout

	out.WebServer.Addr = conf.DashboardAddr
	out.WebServer.Port = conf.DashboardPort
	out.WebServer.User = conf.DashboardUser
	out.WebServer.Password = conf.DashboardPwd
	out.WebServer.AssetsDir = conf.AssetsDir
	if conf.DashboardTLSMode {
		out.WebServer.TLS = &v1.TLSConfig{}
		out.WebServer.TLS.CertFile = conf.DashboardTLSCertFile
		out.WebServer.TLS.KeyFile = conf.DashboardTLSKeyFile
		out.WebServer.PprofEnable = conf.PprofEnabled
	}

	out.Log.To = conf.LogFile
	out.Log.Level = conf.LogLevel
	out.Log.MaxDays = conf.LogMaxDays
	out.Log.DisabledPrintColor = conf.DisableLogColor

	out.DetailedErrorsToClient = lo.ToPtr(conf.DetailedErrorsToClient)
	out.SubDomainHost = conf.SubDomainHost
	out.Custom404Page = conf.Custom404Page
	out.UserConnTimeout = conf.UserConnTimeout
	out.UDPPacketSize = conf.UDPPacketSize
	out.NatHoleAnalysisDataReserveHours = conf.NatHoleAnalysisDataReserveHours

	out.Transport.TCPMux = lo.ToPtr(conf.TCPMux)
	out.Transport.TCPMuxKeepaliveInternal = conf.TCPMuxKeepaliveInterval
	out.Transport.TCPKeepAlive = conf.TCPKeepAlive
	out.Transport.MaxPoolCount = conf.MaxPoolCount
	out.Transport.HeartbeatTimeout = conf.HeartbeatTimeout

	out.Transport.TLS.Force = conf.TLSOnly
	out.Transport.TLS.CertFile = conf.TLSCertFile
	out.Transport.TLS.KeyFile = conf.TLSKeyFile
	out.Transport.TLS.TrustedCaFile = conf.TLSTrustedCaFile

	out.MaxPortsClient = conf.MaxPortsPerClient

	for _, v := range conf.HTTPPlugins {
		out.HTTPPlugins = append(out.HTTPPlugins, v1.HTTPPluginOptions{
			Name:      v.Name,
			Addr:      v.Addr,
			Path:      v.Path,
			Ops:       v.Ops,
			TLSVerify: v.TLSVerify,
		})
	}

	out.AllowPorts, _ = types.NewPortsRangeSliceFromString(conf.AllowPortsStr)
	return out
}
