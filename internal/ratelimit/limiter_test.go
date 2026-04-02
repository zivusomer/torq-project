package ratelimit_test

import (
	"testing"

	"zivusomer/torq-project/internal/ratelimit"
)

func TestLimiterBlocksAfterLimit(t *testing.T) {
	if err := ratelimit.Init(2); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	if !ratelimit.Allow() {
		t.Fatalf("first request should be allowed")
	}
	if !ratelimit.Allow() {
		t.Fatalf("second request should be allowed")
	}
	if ratelimit.Allow() {
		t.Fatalf("third request should be blocked")
	}
}

func TestLimiterRejectsInvalidLimit(t *testing.T) {
	err := ratelimit.Init(0)
	if err == nil {
		t.Fatalf("expected error for invalid limit")
	}
}
