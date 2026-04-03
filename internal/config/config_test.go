package config

import "testing"

func TestLoadFromEnv_Defaults(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_NAME", "")
	t.Setenv("PORT", "")
	t.Setenv("DATASTORE_TYPE", "")
	t.Setenv("DATASTORE_PATH", "")
	t.Setenv("REQUESTS_PER_SECOND", "")
	t.Setenv("RATE_LIMIT_BACKEND", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "")
	t.Setenv("REDIS_KEY_PREFIX", "")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error: %v", err)
	}

	if cfg.AppName != "torq-project-dev" {
		t.Fatalf("expected development file app name, got %q", cfg.AppName)
	}
	if cfg.Env != "development" {
		t.Fatalf("expected default env, got %q", cfg.Env)
	}
	if cfg.Port == "" {
		t.Fatalf("expected non-empty port")
	}
	if cfg.DatastoreType == "" {
		t.Fatalf("expected non-empty datastore type")
	}
	if cfg.DatastorePath == "" {
		t.Fatalf("expected non-empty datastore path")
	}
	if cfg.RequestsPerSecond <= 0 {
		t.Fatalf("expected positive requests per second, got %d", cfg.RequestsPerSecond)
	}
	if cfg.RateLimitBackend == "" {
		t.Fatalf("expected non-empty rate limit backend")
	}
	if cfg.RedisAddr == "" {
		t.Fatalf("expected non-empty redis addr")
	}
	if cfg.RedisDB < 0 {
		t.Fatalf("expected non-negative redis db, got %d", cfg.RedisDB)
	}
	if cfg.RedisKeyPrefix == "" {
		t.Fatalf("expected non-empty redis key prefix")
	}
}

func TestLoadFromEnv_Overrides(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("APP_NAME", "my-app")
	t.Setenv("PORT", "9000")
	t.Setenv("DATASTORE_TYPE", "csv")
	t.Setenv("DATASTORE_PATH", "fixtures/ip.csv")
	t.Setenv("REQUESTS_PER_SECOND", "25")
	t.Setenv("RATE_LIMIT_BACKEND", "redis")
	t.Setenv("REDIS_ADDR", "127.0.0.1:6380")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "4")
	t.Setenv("REDIS_KEY_PREFIX", "my-prefix")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error: %v", err)
	}

	if cfg.AppName != "my-app" {
		t.Fatalf("expected APP_NAME override, got %q", cfg.AppName)
	}
	if cfg.Env != "production" {
		t.Fatalf("expected APP_ENV override, got %q", cfg.Env)
	}
	if cfg.Port != "9000" {
		t.Fatalf("expected PORT override, got %q", cfg.Port)
	}
	if cfg.DatastoreType != "csv" {
		t.Fatalf("expected DATASTORE_TYPE override, got %q", cfg.DatastoreType)
	}
	if cfg.DatastorePath != "fixtures/ip.csv" {
		t.Fatalf("expected DATASTORE_PATH override, got %q", cfg.DatastorePath)
	}
	if cfg.RequestsPerSecond != 25 {
		t.Fatalf("expected REQUESTS_PER_SECOND override, got %d", cfg.RequestsPerSecond)
	}
	if cfg.RateLimitBackend != "redis" {
		t.Fatalf("expected RATE_LIMIT_BACKEND override, got %q", cfg.RateLimitBackend)
	}
	if cfg.RedisAddr != "127.0.0.1:6380" {
		t.Fatalf("expected REDIS_ADDR override, got %q", cfg.RedisAddr)
	}
	if cfg.RedisPassword != "secret" {
		t.Fatalf("expected REDIS_PASSWORD override, got %q", cfg.RedisPassword)
	}
	if cfg.RedisDB != 4 {
		t.Fatalf("expected REDIS_DB override, got %d", cfg.RedisDB)
	}
	if cfg.RedisKeyPrefix != "my-prefix" {
		t.Fatalf("expected REDIS_KEY_PREFIX override, got %q", cfg.RedisKeyPrefix)
	}
}

func TestLoadFromEnv_UsesProductionPresetValues(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("APP_NAME", "")
	t.Setenv("PORT", "")
	t.Setenv("DATASTORE_TYPE", "")
	t.Setenv("DATASTORE_PATH", "")
	t.Setenv("REQUESTS_PER_SECOND", "")
	t.Setenv("RATE_LIMIT_BACKEND", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "")
	t.Setenv("REDIS_KEY_PREFIX", "")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error: %v", err)
	}

	if cfg.AppName == "" {
		t.Fatalf("expected non-empty app name")
	}
	if cfg.Port == "" {
		t.Fatalf("expected non-empty port")
	}
	if cfg.DatastoreType == "" {
		t.Fatalf("expected non-empty datastore type")
	}
	if cfg.DatastorePath == "" {
		t.Fatalf("expected non-empty datastore path")
	}
	if cfg.RequestsPerSecond <= 0 {
		t.Fatalf("expected positive requests per second, got %d", cfg.RequestsPerSecond)
	}
	if cfg.RateLimitBackend == "" {
		t.Fatalf("expected non-empty rate limit backend")
	}
	if cfg.RedisAddr == "" {
		t.Fatalf("expected non-empty redis addr")
	}
	if cfg.RedisDB < 0 {
		t.Fatalf("expected non-negative redis db, got %d", cfg.RedisDB)
	}
	if cfg.RedisKeyPrefix == "" {
		t.Fatalf("expected non-empty redis key prefix")
	}
}

func TestLoadFromEnv_UnsupportedEnv(t *testing.T) {
	t.Setenv("APP_ENV", "staging")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatalf("expected error for unsupported APP_ENV")
	}
}
