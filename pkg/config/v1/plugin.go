package v1

type ClientPluginOptions struct {
}

type TypedClientPluginOptions struct {
	Type string `json:"type"`
	ClientPluginOptions
}
