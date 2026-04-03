package ratelimit

import (
	"sync"
	"time"
)

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
