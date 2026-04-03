package types

import "time"

type LocalBucketState struct {
	Tokens     float64
	LastRefill time.Time
}
