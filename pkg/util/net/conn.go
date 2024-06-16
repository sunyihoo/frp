package net

import (
	"context"
	"net"
)

type ContextGetter interface {
	Context() context.Context
}

type ContextSetter interface {
	WithContext(ctx context.Context)
}

// ContextConn is the connection with context. ContextConn 是有上下文的连接
type ContextConn struct {
	net.Conn

	ctx context.Context
}

func NewContextFromConn(conn net.Conn) context.Context {
	if c, ok := conn.(ContextGetter); ok {
		return c.Context()
	}
	return context.Background()
}

func NewContextConn(ctx context.Context, c net.Conn) *ContextConn {
	return &ContextConn{
		Conn: c,
		ctx:  ctx,
	}
}
