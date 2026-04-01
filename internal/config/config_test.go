package config

import "testing"

func TestLoadFromEnv_Defaults(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_NAME", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("PORT", "")
	t.Setenv("DATASTORE_TYPE", "")
	t.Setenv("DATASTORE_PATH", "")
	t.Setenv("REQUESTS_PER_SECOND", "")

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
	if cfg.LogLevel == "" {
		t.Fatalf("expected non-empty log level")
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
}

func TestLoadFromEnv_Overrides(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("APP_NAME", "my-app")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("PORT", "9000")
	t.Setenv("DATASTORE_TYPE", "csv")
	t.Setenv("DATASTORE_PATH", "fixtures/ip.csv")
	t.Setenv("REQUESTS_PER_SECOND", "25")

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
	if cfg.LogLevel != "debug" {
		t.Fatalf("expected LOG_LEVEL override, got %q", cfg.LogLevel)
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
}

func TestLoadFromEnv_UsesProductionPresetValues(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("APP_NAME", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("PORT", "")
	t.Setenv("DATASTORE_TYPE", "")
	t.Setenv("DATASTORE_PATH", "")
	t.Setenv("REQUESTS_PER_SECOND", "")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error: %v", err)
	}

	if cfg.AppName == "" {
		t.Fatalf("expected non-empty app name")
	}
	if cfg.LogLevel == "" {
		t.Fatalf("expected non-empty log level")
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
}

func TestLoadFromEnv_UnsupportedEnv(t *testing.T) {
	t.Setenv("APP_ENV", "staging")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatalf("expected error for unsupported APP_ENV")
	}
}
