package ratelimit

type Decision struct {
	Allowed           bool
	Limit             int
	Remaining         int
	ResetSeconds      int
	RetryAfterSeconds int
}
