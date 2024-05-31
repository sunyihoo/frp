package nathole

import (
	"sync"
	"time"
)

type RecommandBehavior struct {
	Role              string
	TTL               int
	SendDelayMs       int
	PortsRangeNumber  int
	PortsRandomNumber int
	ListenRandomPorts int
}

type MakeHoleRecords struct {
	mu             sync.Mutex
	scores         []*BehaviorScore
	LastUpdateTime time.Time
}

type BehaviorScore struct {
	Mode  int
	Index int
	// 在-10和10之间
	Score int
}

type Analyzer struct {
	// 键名是客户端IP+访客IP
	records             map[string]*MakeHoleRecords
	dataReserveDuration time.Duration

	my sync.Mutex
}
