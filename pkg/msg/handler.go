package msg

import (
	"io"
	"reflect"
)

// Dispatcher 用于向 net.Conn 发送消息或寄存器处理(register handles)程序，用于从 net.Conn 读取的消息。
type Dispatcher struct {
	rw io.ReadWriter

	sendCh         chan Message
	doneCh         chan struct{}
	msgHandlers    map[reflect.Type]func(Message)
	defaultHandler func(Message)
}
