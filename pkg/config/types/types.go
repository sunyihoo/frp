package types

type PortsRange struct {
	Start  int `json:"start,omitempty"`
	End    int `json:"end,omitempty"`
	Single int `json:"single,omitempty"`
}
