package net

import (
	"context"
	"net"
)

// ContextConn is the connection with context. ContextConn 是有上下文的连接
type ContextConn struct {
	net.Conn

	ctx context.Context
}

func NewContextConn(ctx context.Context, c net.Conn) *ContextConn {
	return &ContextConn{
		Conn: c,
		ctx:  ctx,
	}
}
