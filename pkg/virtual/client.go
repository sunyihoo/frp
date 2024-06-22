package virtual

import (
	"github.com/sunyihoo/frp/client"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
)

type Client struct {
	l   *netpkg.InternalListener
	svr *client.Service
}
