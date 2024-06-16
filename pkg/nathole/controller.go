package nathole

import (
	"github.com/sunyihoo/frp/pkg/msg"
	"github.com/sunyihoo/frp/pkg/transport"
	"sync"
	"time"
)

type ClientCfg struct {
	name       string
	sk         string
	allowUsers []string
	sidCh      chan string
}

type Session struct {
	sid            string
	analysisKey    string
	recommandMode  int
	recommandIndex int

	visitorMsg         *msg.NatHoleVisitor
	visitorTransporter transport.MessageTransporter
	vResp              *msg.NatHoleResp
	vNatFeature        *NatFeature
	vBehavior          RecommandBehavior

	clientMsg         *msg.NatHoleClient
	clientTransporter transport.MessageTransporter
	cResp             *msg.NatHoleResp
	cNatFeature       *NatFeature
	cBehavior         RecommandBehavior

	notifyCh chan struct{}
}

type Controller struct {
	clientCfgs map[string]*ClientCfg
	sessions   map[string]*Session
	analyzer   *Analyzer

	mu sync.Mutex
}

func NewController(analysisDataReserveDuration time.Duration) (*Controller, error) {
	return &Controller{
		clientCfgs: make(map[string]*ClientCfg),
		sessions:   make(map[string]*Session),
		analyzer:   NewAnalyzer(analysisDataReserveDuration),
	}, nil
}
