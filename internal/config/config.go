package config

import "os"

type Config struct {
	AppName  string
	Env      string
	LogLevel string
}

func LoadFromEnv() Config {
	return Config{
		AppName:  getEnv("APP_NAME", "torq-project"),
		Env:      getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
