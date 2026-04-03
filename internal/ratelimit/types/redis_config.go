package types

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
}
