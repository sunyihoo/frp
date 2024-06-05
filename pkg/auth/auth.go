package auth

import "github.com/sunyihoo/frp/pkg/msg"

type Verifier interface {
	VerifyLogin(*msg.Login) error
	VerifyPing(*msg.Ping) error
	VerifyNewWorkConn(conn *msg.NewWorkConn) error
}
