package ratelimit

import (
	"fmt"
	"math"
	"sync"
	"time"
)

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
	bucket := l.getOrCreateBucketForKey(key, now)
	l.refillBucketTokens(bucket, now)

	if bucket.tokens >= 1 {
		remaining := l.consumeToken(bucket)
		return l.allowedDecision(remaining)
	}

	retryAfter := l.retryAfterSeconds(bucket.tokens)
	return l.deniedDecision(bucket.tokens, retryAfter)
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

func (l *Limiter) consumeToken(bucket *bucketState) int {
	bucket.tokens--
	return int(math.Floor(bucket.tokens))
}

func (l *Limiter) getOrCreateBucketForKey(key string, now time.Time) *bucketState {
	bucket, ok := l.buckets[key]
	if ok {
		return bucket
	}
	return l.createBucketForKey(key, now)
}

func (l *Limiter) refillBucketTokens(bucket *bucketState, now time.Time) {
	elapsed := l.elapsedSeconds(now, bucket.lastRefill)
	if elapsed > 0 {
		bucket.tokens = math.Min(l.capacity, bucket.tokens+(elapsed*l.refillRate))
		bucket.lastRefill = now
	}
}

func (l *Limiter) elapsedSeconds(now time.Time, lastRefill time.Time) float64 {
	return now.Sub(lastRefill).Seconds()
}

func (l *Limiter) createBucketForKey(key string, now time.Time) *bucketState {
	bucket := &bucketState{
		tokens:     l.capacity,
		lastRefill: now,
	}
	l.buckets[key] = bucket
	return bucket
}

func (l *Limiter) retryAfterSeconds(tokens float64) int {
	retryAfter := int(math.Ceil((1 - tokens) / l.refillRate))
	if retryAfter < 1 {
		return 1
	}
	return retryAfter
}

func (l *Limiter) resetSeconds(tokens float64) int {
	reset := int(math.Ceil((l.capacity - tokens) / l.refillRate))
	if reset < 0 {
		return 0
	}
	return reset
}

func (l *Limiter) allowedDecision(remaining int) Decision {
	return Decision{
		Allowed:      true,
		Limit:        l.limit,
		Remaining:    remaining,
		ResetSeconds: 0,
	}
}

func (l *Limiter) deniedDecision(tokens float64, retryAfter int) Decision {
	remaining := int(math.Floor(math.Max(tokens-1, 0)))
	reset := l.resetSeconds(tokens)
	if reset < retryAfter {
		reset = retryAfter
	}

	return Decision{
		Allowed:           false,
		Limit:             l.limit,
		Remaining:         remaining,
		ResetSeconds:      reset,
		RetryAfterSeconds: retryAfter,
	}
}
