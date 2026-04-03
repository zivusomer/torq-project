package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"zivusomer/torq-project/internal/ratelimit/types"
)

func NewBackend(limit int, cfg types.RedisConfig) (types.Backend, error) {
	if limit < 1 {
		return nil, fmt.Errorf("rate limit must be >= 1")
	}
	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis addr is required")
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "ratelimit"
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	l := &types.RedisLimiter{
		Client:      client,
		Limit:       limit,
		Capacity:    float64(limit),
		RefillPerMs: float64(limit) / 1000.0,
		KeyPrefix:   cfg.Prefix,
	}
	l.AllowForKeyFn = func(key string) types.Decision {
		return allowForKey(l, key)
	}
	return l, nil
}

func allowForKey(l *types.RedisLimiter, key string) types.Decision {
	if key == "" {
		key = "anonymous"
	}

	nowMs := time.Now().UnixMilli()
	redisKey := l.KeyPrefix + ":" + key
	ttlSeconds := strconv.Itoa(3600)

	result, err := allowScript.Run(
		context.Background(),
		l.Client,
		[]string{redisKey},
		nowMs,
		l.Capacity,
		l.RefillPerMs,
		ttlSeconds,
	).Result()
	if err != nil {
		return types.Decision{Allowed: false, Limit: l.Limit, Remaining: 0, ResetSeconds: 1, RetryAfterSeconds: 1}
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 4 {
		return types.Decision{Allowed: false, Limit: l.Limit, Remaining: 0, ResetSeconds: 1, RetryAfterSeconds: 1}
	}

	allowed := values[0].(int64) == 1
	remaining := int(values[1].(int64))
	reset := int(values[2].(int64))
	retry := int(values[3].(int64))

	return types.Decision{
		Allowed:           allowed,
		Limit:             l.Limit,
		Remaining:         remaining,
		ResetSeconds:      reset,
		RetryAfterSeconds: retry,
	}
}
