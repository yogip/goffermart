package helpers

import (
	"sync"
	"time"
)

type StopLimiter struct {
	lock      *sync.RWMutex
	stopUnill time.Time
}

func NewStopLimiter() *StopLimiter {
	return &StopLimiter{
		lock:      &sync.RWMutex{},
		stopUnill: time.Now(),
	}
}

func (s *StopLimiter) CanWork() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return time.Now().After(s.stopUnill)
}

func (s *StopLimiter) EnsureLimit() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	sleep := time.Until(s.stopUnill)
	time.Sleep(sleep)
}

func (s *StopLimiter) StopUntil(interval time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stopUnill = time.Now().Add(interval)
}
