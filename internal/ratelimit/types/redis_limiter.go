package types

import goredis "github.com/redis/go-redis/v9"

type RedisLimiter struct {
	Client        *goredis.Client
	Limit         int
	Capacity      float64
	RefillPerMs   float64
	KeyPrefix     string
	AllowForKeyFn func(string) Decision
}

func (l *RedisLimiter) AllowForKey(key string) Decision {
	if l.AllowForKeyFn == nil {
		return Decision{Allowed: false, Limit: l.Limit, Remaining: 0, ResetSeconds: 1, RetryAfterSeconds: 1}
	}
	return l.AllowForKeyFn(key)
}
