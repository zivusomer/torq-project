package local

import (
	"fmt"
	"math"
	"time"

	"zivusomer/torq-project/internal/ratelimit/types"
)

func NewBackend(limit int) (types.Backend, error) {
	return newLimiter(limit)
}

func AllowForKey(l *types.LocalLimiter, key string) types.Decision {
	return allow(l, key)
}

func allow(l *types.LocalLimiter, key string) types.Decision {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	if key == "" {
		key = "anonymous"
	}

	now := time.Now()
	bucket := getOrCreateBucketForKey(l, key, now)
	refillBucketTokens(l, bucket, now)

	if bucket.Tokens >= 1 {
		remaining := consumeToken(bucket)
		return allowedDecision(l, remaining)
	}

	retryAfter := retryAfterSeconds(l, bucket.Tokens)
	return deniedDecision(l, bucket.Tokens, retryAfter)
}

func newLimiter(limit int) (*types.LocalLimiter, error) {
	if limit < 1 {
		return nil, fmt.Errorf("rate limit must be >= 1")
	}
	l := &types.LocalLimiter{
		Limit:      limit,
		Capacity:   float64(limit),
		RefillRate: float64(limit),
		Buckets:    make(map[string]*types.LocalBucketState),
	}
	l.AllowForKeyFn = func(key string) types.Decision {
		return allow(l, key)
	}
	return l, nil
}

func consumeToken(bucket *types.LocalBucketState) int {
	bucket.Tokens--
	return int(math.Floor(bucket.Tokens))
}

func getOrCreateBucketForKey(l *types.LocalLimiter, key string, now time.Time) *types.LocalBucketState {
	bucket, ok := l.Buckets[key]
	if ok {
		return bucket
	}
	return createBucketForKey(l, key, now)
}

func refillBucketTokens(l *types.LocalLimiter, bucket *types.LocalBucketState, now time.Time) {
	elapsed := elapsedSeconds(now, bucket.LastRefill)
	if elapsed > 0 {
		bucket.Tokens = math.Min(l.Capacity, bucket.Tokens+(elapsed*l.RefillRate))
		bucket.LastRefill = now
	}
}

func elapsedSeconds(now time.Time, lastRefill time.Time) float64 {
	return now.Sub(lastRefill).Seconds()
}

func createBucketForKey(l *types.LocalLimiter, key string, now time.Time) *types.LocalBucketState {
	bucket := &types.LocalBucketState{
		Tokens:     l.Capacity,
		LastRefill: now,
	}
	l.Buckets[key] = bucket
	return bucket
}

func retryAfterSeconds(l *types.LocalLimiter, tokens float64) int {
	retryAfter := int(math.Ceil((1 - tokens) / l.RefillRate))
	if retryAfter < 1 {
		return 1
	}
	return retryAfter
}

func resetSeconds(l *types.LocalLimiter, tokens float64) int {
	reset := int(math.Ceil((l.Capacity - tokens) / l.RefillRate))
	if reset < 0 {
		return 0
	}
	return reset
}

func allowedDecision(l *types.LocalLimiter, remaining int) types.Decision {
	return types.Decision{
		Allowed:      true,
		Limit:        l.Limit,
		Remaining:    remaining,
		ResetSeconds: 0,
	}
}

func deniedDecision(l *types.LocalLimiter, tokens float64, retryAfter int) types.Decision {
	remaining := int(math.Floor(math.Max(tokens-1, 0)))
	reset := resetSeconds(l, tokens)
	if reset < retryAfter {
		reset = retryAfter
	}

	return types.Decision{
		Allowed:           false,
		Limit:             l.Limit,
		Remaining:         remaining,
		ResetSeconds:      reset,
		RetryAfterSeconds: retryAfter,
	}
}
