package client

// Connector 是用于建立与服务器连接的接口。
type Connector interface {
	Open() error
	Connect() error
	Close() error
}
