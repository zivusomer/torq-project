package ratelimit_test

import (
	"testing"

	"zivusomer/torq-project/internal/ratelimit"
)

func TestLimiterBlocksAfterLimit(t *testing.T) {
	if err := ratelimit.Init(2); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	if !ratelimit.AllowForKey("user-1").Allowed {
		t.Fatalf("first request should be allowed")
	}
	if !ratelimit.AllowForKey("user-1").Allowed {
		t.Fatalf("second request should be allowed")
	}
	decision := ratelimit.AllowForKey("user-1")
	if decision.Allowed {
		t.Fatalf("third request should be blocked")
	}
	if decision.RetryAfterSeconds < 1 {
		t.Fatalf("expected retry-after to be positive, got %d", decision.RetryAfterSeconds)
	}
}

func TestLimiterRejectsInvalidLimit(t *testing.T) {
	err := ratelimit.Init(0)
	if err == nil {
		t.Fatalf("expected error for invalid limit")
	}
}

func TestLimiterIsolatedByKey(t *testing.T) {
	if err := ratelimit.Init(1); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	if !ratelimit.AllowForKey("user-a").Allowed {
		t.Fatalf("first request for user-a should be allowed")
	}
	if !ratelimit.AllowForKey("user-b").Allowed {
		t.Fatalf("first request for user-b should be allowed")
	}
}
