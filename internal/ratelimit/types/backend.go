package types

type Backend interface {
	AllowForKey(key string) Decision
}
