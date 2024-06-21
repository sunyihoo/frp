package net

import (
	"context"
	"errors"
	"github.com/fatedier/golib/crypto"
	"github.com/quic-go/quic-go"
	"github.com/sunyihoo/frp/pkg/util/xlog"
	"io"
	"net"
	"time"
)

type ContextGetter interface {
	Context() context.Context
}

type ContextSetter interface {
	WithContext(ctx context.Context)
}

func NewLogFromConn(conn net.Conn) *xlog.Logger {
	if c, ok := conn.(ContextGetter); ok {
		return xlog.FromContextSafe(c.Context())
	}
	return xlog.New()
}

func NewContextFromConn(conn net.Conn) context.Context {
	if c, ok := conn.(ContextGetter); ok {
		return c.Context()
	}
	return context.Background()
}

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

type WrapReadWriteCloserConn struct {
	io.ReadWriteCloser

	underConn net.Conn

	remoteAddr net.Addr
}

func WrapReadWriteCloserToConn(rwc io.ReadWriteCloser, underConn net.Conn) *WrapReadWriteCloserConn {
	return &WrapReadWriteCloserConn{
		ReadWriteCloser: rwc,
		underConn:       underConn,
	}
}

func (conn *WrapReadWriteCloserConn) LocalAddr() net.Addr {
	if conn.underConn != nil {
		return conn.underConn.LocalAddr()
	}
	// todo 学习
	return (*net.TCPAddr)(nil)
}

func (conn *WrapReadWriteCloserConn) SetRemoteAddr(addr net.Addr) {
	conn.remoteAddr = addr
}

func (conn *WrapReadWriteCloserConn) RemoteAddr() net.Addr {
	if conn.remoteAddr != nil {
		return conn.remoteAddr
	}
	if conn.underConn != nil {
		return conn.underConn.RemoteAddr()
	}
	return (*net.TCPAddr)(nil)
}

func (conn *WrapReadWriteCloserConn) SetDeadline(t time.Time) error {
	if conn.underConn != nil {
		return conn.underConn.SetDeadline(t)
	}
	// todo 学习
	return &net.OpError{Op: "set", Net: "wrap", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (conn *WrapReadWriteCloserConn) SetReadDeadline(t time.Time) error {
	if conn.underConn != nil {
		return conn.underConn.SetReadDeadline(t)
	}
	return &net.OpError{Op: "set", Net: "wrap", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (conn *WrapReadWriteCloserConn) SetWriteDeadline(t time.Time) error {
	if conn.underConn != nil {
		return conn.underConn.SetWriteDeadline(t)
	}
	return &net.OpError{Op: "set", Net: "wrap", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

type wrapQuicStream struct {
	quic.Stream
	c quic.Connection
}

func QuicStreamToNetConn(s quic.Stream, c quic.Connection) net.Conn {
	return &wrapQuicStream{
		Stream: s,
		c:      c,
	}
}

func (conn *wrapQuicStream) LocalAddr() net.Addr {
	if conn.c != nil {
		return conn.c.LocalAddr()
	}
	return (*net.TCPAddr)(nil)
}

func (conn *wrapQuicStream) RemoteAddr() net.Addr {
	if conn.c != nil {
		return conn.c.RemoteAddr()
	}
	return (*net.TCPAddr)(nil)
}

func (conn *wrapQuicStream) Close() error {
	conn.Stream.CancelRead(0)
	return conn.Stream.Close()
}

func NewCryptoReadWriter(rw io.ReadWriter, key []byte) (io.ReadWriter, error) {
	encReader := crypto.NewReader(rw, key)
	encWriter, err := crypto.NewWriter(rw, key)
	if err != nil {
		return nil, err
	}
	return struct {
		io.Reader
		io.Writer
	}{
		Reader: encReader,
		Writer: encWriter,
	}, nil
}
