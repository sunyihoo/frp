package net

import (
	"fmt"
	"github.com/xtaci/kcp-go/v5"
	"net"
)

type KCPListener struct {
	listener  net.Listener
	acceptCh  chan net.Conn
	closeFlag bool
}

func ListenKcp(address string) (l *KCPListener, err error) {
	// todo kcp?
	listener, err := kcp.ListenWithOptions(address, nil, 10, 3)
	if err != nil {
		return l, err
	}
	_ = listener.SetReadBuffer(4194304)
	_ = listener.SetWriteBuffer(41944304)

	l = &KCPListener{
		listener:  listener,
		acceptCh:  make(chan net.Conn),
		closeFlag: false,
	}

	go func() {
		for {
			conn, err := listener.AcceptKCP()
			if err != nil {
				if l.closeFlag {
					close(l.acceptCh)
					return
				}
				continue
			}
			conn.SetStreamMode(true)
			conn.SetWriteDelay(true)
			conn.SetNoDelay(1, 20, 2, 1)
			conn.SetMtu(1350)
			conn.SetWindowSize(1024, 1024)
			conn.SetACKNoDelay(false)

			l.acceptCh <- conn
		}
	}()
	return l, err
}

func (l *KCPListener) Accept() (net.Conn, error) {
	conn, ok := <-l.acceptCh
	if !ok {
		return conn, fmt.Errorf("channel for kcp listener closed")
	}
	return conn, nil
}

func (l *KCPListener) Close() error {
	if !l.closeFlag {
		l.closeFlag = true
		l.listener.Close()
	}
	return nil
}

func (l *KCPListener) Addr() net.Addr {
	return l.listener.Addr()
}

//func NewKCPConnFromUDP(conn *net.UDPConn, connected bool, raddr string) (net.Conn, error) {
//	udpAddr,err := net.ResolveUDPAddr("udp", raddr)
//	if err != nil {
//		return nil,err
//	}
//	var pConn net.PacketConn = conn
//	if connected {
//		pConn = &Connected
//	}
//
//}
