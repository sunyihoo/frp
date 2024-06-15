package ssh

import (
	"fmt"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/transport"
	"github.com/sunyihoo/frp/pkg/util/log"
	netpkg "github.com/sunyihoo/frp/pkg/util/net"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"strconv"
	"strings"
)

type Gateway struct {
	bindPort int
	ln       net.Listener

	peerServerListener *netpkg.InternalListener

	sshConfig *ssh.ServerConfig
}

func NewGateway(
	cfg v1.SSHTunnelGateway, bindAddr string,
	peerServerListener *netpkg.InternalListener,
) (*Gateway, error) {
	sshConfig := &ssh.ServerConfig{}

	// privateKey
	var (
		privateKeyBytes []byte
		err             error
	)
	if cfg.PrivateKeyFile != "" {
		privateKeyBytes, err = os.ReadFile(cfg.PrivateKeyFile)
	} else {
		if cfg.AutoGenPrivateKeyPath != "" {
			privateKeyBytes, _ = os.ReadFile(cfg.AutoGenPrivateKeyPath)
		}
		if len(privateKeyBytes) == 0 {
			privateKeyBytes, err = transport.NewRandomPrivateKey()
			if err == nil && cfg.AutoGenPrivateKeyPath != "" {
				err = os.WriteFile(cfg.AutoGenPrivateKeyPath, privateKeyBytes, 0o600)
			}
		}
	}
	if err != nil {
		return nil, err
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	sshConfig.AddHostKey(privateKey)

	sshConfig.NoClientAuth = cfg.AuthorizedKeysFile == ""
	sshConfig.PublicKeyCallback = func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		authorizedKeysMap, err := loadAuthorizedKeysFromFile(cfg.AuthorizedKeysFile)
		if err != nil {
			log.Errorf("load authorized keys file error: %v", err)
			return nil, fmt.Errorf("internal error")
		}

		user, ok := authorizedKeysMap[string(key.Marshal())]
		if !ok {
			return nil, fmt.Errorf("unknown public key for remoteAddr %q", conn.RemoteAddr())
		}
		return &ssh.Permissions{
			Extensions: map[string]string{
				"user": user,
			},
		}, nil
	}

	ln, err := net.Listen("tcp", net.JoinHostPort(bindAddr, strconv.Itoa(cfg.BindPort)))
	if err != nil {
		return nil, err
	}
	return &Gateway{
		bindPort:           cfg.BindPort,
		ln:                 ln,
		peerServerListener: peerServerListener,
		sshConfig:          sshConfig,
	}, nil
}

func loadAuthorizedKeysFromFile(path string) (map[string]string, error) {
	authorizedKeysMap := make(map[string]string) // value is username
	authorizedKeysBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	for len(authorizedKeysBytes) > 0 {
		pubKey, comment, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return nil, err
		}

		authorizedKeysMap[string(pubKey.Marshal())] = strings.TrimSpace(comment)
		authorizedKeysBytes = rest
	}
	return authorizedKeysMap, nil
}
