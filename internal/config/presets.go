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
	},
	"production": {
		AppName:           "torq-project",
		Env:               "production",
		Port:              "8080",
		DatastoreType:     "csv",
		DatastorePath:     "data/ip_locations.csv",
		RequestsPerSecond: 3,
	},
}

func presetForEnv(env string) (Config, error) {
	cfg, ok := presets[env]
	if !ok {
		return Config{}, fmt.Errorf("unsupported APP_ENV: %s", env)
	}
	return cfg, nil
}
