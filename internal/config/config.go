package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppName           string
	Env               string
	Port              string
	DatastoreType     string
	DatastorePath     string
	RequestsPerSecond int
}

func LoadFromEnv() (Config, error) {
	env := getEnv("APP_ENV", "development")
	cfg, err := presetForEnv(env)
	if err != nil {
		return Config{}, err
	}

	// Environment variables have highest precedence over environment presets.
	cfg.AppName = getEnv("APP_NAME", cfg.AppName)
	cfg.Env = env
	cfg.Port = getEnv("PORT", cfg.Port)
	cfg.DatastoreType = getEnv("DATASTORE_TYPE", cfg.DatastoreType)
	cfg.DatastorePath = getEnv("DATASTORE_PATH", cfg.DatastorePath)
	cfg.RequestsPerSecond = getEnvInt("REQUESTS_PER_SECOND", cfg.RequestsPerSecond)

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}

	return parsed
}
