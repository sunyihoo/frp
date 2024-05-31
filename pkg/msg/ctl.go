package msg

import jsonMsg "github.com/fatedier/golib/msg/json"

type Message = jsonMsg.Message

var msgCtl *jsonMsg.MsgCtl

func init() {
	msgCtl = jsonMsg.NewMsgCtl()
	for typeByte, msg := range msgTypeMap {
		msgCtl.RegisterMsg(typeByte, msg)
	}
}
