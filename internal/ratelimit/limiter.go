package ratelimit

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Decision struct {
	Allowed           bool
	Limit             int
	Remaining         int
	ResetSeconds      int
	RetryAfterSeconds int
}

type Limiter struct {
	mu         sync.Mutex
	limit      int
	capacity   float64
	refillRate float64
	buckets    map[string]*bucketState
}

type bucketState struct {
	tokens     float64
	lastRefill time.Time
}

var (
	globalMu      sync.RWMutex
	globalLimiter *Limiter
)

func Init(limit int) error {
	limiter, err := newLimiter(limit)
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
	return limiter.allow("global").Allowed
}

func AllowForKey(key string) Decision {
	globalMu.RLock()
	limiter := globalLimiter
	globalMu.RUnlock()
	if limiter == nil {
		return Decision{Allowed: false}
	}
	return limiter.allow(key)
}

func (l *Limiter) allow(key string) Decision {
	l.mu.Lock()
	defer l.mu.Unlock()

	if key == "" {
		key = "anonymous"
	}

	now := time.Now()
	bucket, ok := l.buckets[key]
	if !ok {
		bucket = &bucketState{
			tokens:     l.capacity,
			lastRefill: now,
		}
		l.buckets[key] = bucket
	}

	elapsed := now.Sub(bucket.lastRefill).Seconds()
	if elapsed > 0 {
		bucket.tokens = math.Min(l.capacity, bucket.tokens+(elapsed*l.refillRate))
		bucket.lastRefill = now
	}

	decision := Decision{
		Allowed:      false,
		Limit:        l.limit,
		Remaining:    int(math.Floor(math.Max(bucket.tokens-1, 0))),
		ResetSeconds: int(math.Ceil((l.capacity - bucket.tokens) / l.refillRate)),
	}
	if decision.ResetSeconds < 0 {
		decision.ResetSeconds = 0
	}

	if bucket.tokens >= 1 {
		bucket.tokens--
		decision.Allowed = true
		decision.Remaining = int(math.Floor(bucket.tokens))
		return decision
	}

	retryAfter := int(math.Ceil((1 - bucket.tokens) / l.refillRate))
	if retryAfter < 1 {
		retryAfter = 1
	}
	decision.RetryAfterSeconds = retryAfter
	if decision.ResetSeconds < retryAfter {
		decision.ResetSeconds = retryAfter
	}
	return decision
}

func newLimiter(limit int) (*Limiter, error) {
	if limit < 1 {
		return nil, fmt.Errorf("rate limit must be >= 1")
	}
	return &Limiter{
		limit:      limit,
		capacity:   float64(limit),
		refillRate: float64(limit),
		buckets:    make(map[string]*bucketState),
	}, nil
}
