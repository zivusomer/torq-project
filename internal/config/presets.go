package config

import "fmt"

var presets = map[string]Config{
	"development": {
		AppName:           "torq-project-dev",
		Env:               "development",
		Port:              "8080",
		DatastoreType:     "csv",
		DatastorePath:     "data/ip_locations.csv",
		RequestsPerSecond: 2,
		RateLimitBackend:  "inmemory",
		RedisAddr:         "localhost:6379",
		RedisDB:           0,
		RedisKeyPrefix:    "torq:ratelimit",
	},
	"production": {
		AppName:           "torq-project",
		Env:               "production",
		Port:              "8080",
		DatastoreType:     "csv",
		DatastorePath:     "data/ip_locations.csv",
		RequestsPerSecond: 3,
		RateLimitBackend:  "redis",
		RedisAddr:         "localhost:6379",
		RedisDB:           0,
		RedisKeyPrefix:    "torq:ratelimit",
	},
}

func presetForEnv(env string) (Config, error) {
	cfg, ok := presets[env]
	if !ok {
		return Config{}, fmt.Errorf("unsupported APP_ENV: %s", env)
	}
	return cfg, nil
}
