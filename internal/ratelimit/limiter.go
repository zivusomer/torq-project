package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	limit    int
	window   int64
	requests int
}

func New(limit int) *Limiter {
	if limit < 1 {
		limit = 1
	}
	return &Limiter{limit: limit}
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
