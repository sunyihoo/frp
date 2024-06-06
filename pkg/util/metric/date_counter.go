package metric

import (
	"sync"
	"time"
)

type DateCounter interface {
	TodayCount() int64
	GetLastDaysCount(lastdays int64) []int64
	Inc(int64)
	Dec(int64)
	Snapshot() DateCounter
	Clear()
}

func NewDateCounter(reserveDays int64) DateCounter {
	if reserveDays <= 0 {
		reserveDays = 1
	}
	return newStandardDateCounter(reserveDays)
}

type StandardDateCounter struct {
	reserveDays int64
	counts      []int64

	lastUpdateDate time.Time
	mu             sync.Mutex
}

func newStandardDateCounter(reserveDays int64) *StandardDateCounter {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	s := &StandardDateCounter{
		reserveDays:    reserveDays,
		counts:         make([]int64, reserveDays),
		lastUpdateDate: now,
	}
	return s
}
