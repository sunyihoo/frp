package v1

type ProxyBaseConfig struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Transport
}

type ProxyConfigurer interface {
	Complete() (namePrefix string)
	GetBaseConfig() *Proxy
}
