package config

import "testing"

func TestLoadFromEnv_Defaults(t *testing.T) {
	t.Setenv("APP_NAME", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("LOG_LEVEL", "")

	cfg := LoadFromEnv()

	if cfg.AppName != "torq-project" {
		t.Fatalf("expected default app name, got %q", cfg.AppName)
	}
	if cfg.Env != "development" {
		t.Fatalf("expected default env, got %q", cfg.Env)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("expected default log level, got %q", cfg.LogLevel)
	}
}

func TestLoadFromEnv_Overrides(t *testing.T) {
	t.Setenv("APP_NAME", "my-app")
	t.Setenv("APP_ENV", "production")
	t.Setenv("LOG_LEVEL", "debug")

	cfg := LoadFromEnv()

	if cfg.AppName != "my-app" {
		t.Fatalf("expected APP_NAME override, got %q", cfg.AppName)
	}
	if cfg.Env != "production" {
		t.Fatalf("expected APP_ENV override, got %q", cfg.Env)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("expected LOG_LEVEL override, got %q", cfg.LogLevel)
	}
}
