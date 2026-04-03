package types

import "sync"

type LocalLimiter struct {
	Mu            sync.Mutex
	Limit         int
	Capacity      float64
	RefillRate    float64
	Buckets       map[string]*LocalBucketState
	AllowForKeyFn func(string) Decision
}

func (l *LocalLimiter) AllowForKey(key string) Decision {
	if l.AllowForKeyFn == nil {
		return Decision{Allowed: false, Limit: l.Limit, Remaining: 0, ResetSeconds: 1, RetryAfterSeconds: 1}
	}
	return l.AllowForKeyFn(key)
}
