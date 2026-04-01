package ratelimit

import "testing"

func TestLimiterBlocksAfterLimit(t *testing.T) {
	limiter := New(2)

	if !limiter.Allow() {
		t.Fatalf("first request should be allowed")
	}
	if !limiter.Allow() {
		t.Fatalf("second request should be allowed")
	}
	if limiter.Allow() {
		t.Fatalf("third request should be blocked")
	}
}
