package ratelimit

import (
	"sync"

	"zivusomer/torq-project/internal/ratelimit/local"
	redisbackend "zivusomer/torq-project/internal/ratelimit/redis"
	"zivusomer/torq-project/internal/ratelimit/types"
)

type Backend = types.Backend
type Decision = types.Decision
type RedisConfig = types.RedisConfig

var (
	globalMu      sync.RWMutex
	globalBackend Backend
)

func InitLocal(limit int) error {
	backend, err := local.NewBackend(limit)
	if err != nil {
		return err
	}
	return initBackend(backend)
}

func InitRedis(limit int, cfg RedisConfig) error {
	backend, err := redisbackend.NewBackend(limit, cfg)
	if err != nil {
		return err
	}
	return initBackend(backend)
}

func initBackend(backend Backend) error {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalBackend = backend
	return nil
}

func AllowForKey(key string) Decision {
	globalMu.RLock()
	backend := globalBackend
	globalMu.RUnlock()
	if backend == nil {
		return Decision{Allowed: false}
	}
	return backend.AllowForKey(key)
}
