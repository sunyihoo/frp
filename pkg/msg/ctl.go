package msg

import (
	jsonMsg "github.com/fatedier/golib/msg/json"
	"io"
)

type Message = jsonMsg.Message

var msgCtl *jsonMsg.MsgCtl

func init() {
	msgCtl = jsonMsg.NewMsgCtl()
	for typeByte, msg := range msgTypeMap {
		msgCtl.RegisterMsg(typeByte, msg)
	}
}

func ReadMsg(c io.Reader) (msg Message, err error) {
	return msgCtl.ReadMsg(c)
}
