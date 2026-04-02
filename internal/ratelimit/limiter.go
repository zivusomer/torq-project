package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	limit    int
	window   int64
	requests int
}

var (
	globalMu      sync.RWMutex
	globalLimiter *Limiter
)

func New(limit int) (*Limiter, error) {
	if limit < 1 {
		return nil, fmt.Errorf("rate limit must be >= 1")
	}
	return &Limiter{limit: limit}, nil
}

func Init(limit int) error {
	limiter, err := New(limit)
	if err != nil {
		return err
	}

	globalMu.Lock()
	defer globalMu.Unlock()
	globalLimiter = limiter
	return nil
}

func Allow() bool {
	globalMu.RLock()
	limiter := globalLimiter
	globalMu.RUnlock()
	if limiter == nil {
		return false
	}
	return limiter.Allow()
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Unix()
	if now != l.window {
		l.window = now
		l.requests = 0
	}

	if l.requests >= l.limit {
		return false
	}

	l.requests++
	return true
}
