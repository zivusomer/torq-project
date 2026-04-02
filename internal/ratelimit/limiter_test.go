package ratelimit

import "testing"

func TestLimiterBlocksAfterLimit(t *testing.T) {
	limiter, err := New(2)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

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

func TestLimiterRejectsInvalidLimit(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatalf("expected error for invalid limit")
	}
}
