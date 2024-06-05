package ssh

import (
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"golang.org/x/crypto/ssh"
	"net"
)

type GateWay struct {
	bindPort int
	ln       net.Listener

	peerServerListener *netpkg.InternalListener

	sshConfig *ssh.ServerConfig
}
