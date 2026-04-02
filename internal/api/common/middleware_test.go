package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"zivusomer/torq-project/internal/httpserver/httpx"
	"zivusomer/torq-project/internal/ratelimit"
)

func TestCallerIPIdentityMiddleware_FromForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=1.1.1.1", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 10.0.0.1")
	ctx := ContextWithDefaults(httptest.NewRecorder(), req)

	CallerIPIdentityMiddleware(ctx, func(*Context) {})

	if ctx.CallerKey != "203.0.113.1" {
		t.Fatalf("caller key = %q, want %q", ctx.CallerKey, "203.0.113.1")
	}
}

func TestRequestIPMiddleware_Valid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=8.8.8.8", nil)
	ctx := ContextWithDefaults(httptest.NewRecorder(), req)
	called := false

	RequestIPMiddleware(ctx, func(*Context) {
		called = true
	})

	if !called {
		t.Fatalf("expected next middleware to be called")
	}
	if got := ctx.RequestedIP.String(); got != "8.8.8.8" {
		t.Fatalf("requested ip = %q, want %q", got, "8.8.8.8")
	}
}

func TestRequestIPMiddleware_MissingIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/find-country", nil)
	ctx := ContextWithDefaults(httptest.NewRecorder(), req)

	RequestIPMiddleware(ctx, func(*Context) {
		t.Fatalf("next middleware should not be called")
	})

	if ctx.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", ctx.StatusCode, http.StatusBadRequest)
	}
}

func TestRateLimitMiddleware_RateLimitedIncludesHeaders(t *testing.T) {
	if err := ratelimit.Init(1); err != nil {
		t.Fatalf("ratelimit.Init() error: %v", err)
	}

	allowReq := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=8.8.8.8", nil)
	allowCtx := ContextWithDefaults(httptest.NewRecorder(), allowReq)
	allowCtx.CallerKey = "rate-limit-test"
	RateLimitMiddleware(allowCtx, func(*Context) {})

	limitedReq := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=8.8.8.8", nil)
	limitedCtx := ContextWithDefaults(httptest.NewRecorder(), limitedReq)
	limitedCtx.CallerKey = "rate-limit-test"
	RateLimitMiddleware(limitedCtx, func(*Context) {
		t.Fatalf("next middleware should not be called when rate limited")
	})

	if limitedCtx.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", limitedCtx.StatusCode, http.StatusTooManyRequests)
	}
	if limitedCtx.ResponseHeader.Get("Retry-After") == "" {
		t.Fatalf("expected Retry-After header")
	}
	if limitedCtx.ResponseHeader.Get("RateLimit-Limit") == "" {
		t.Fatalf("expected RateLimit-Limit header")
	}
}

func TestWriteHTTPError_WithHTTPError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=1.1.1.1", nil)
	ctx := ContextWithDefaults(httptest.NewRecorder(), req)

	WriteHTTPError(ctx, &httpx.Error{StatusCode: http.StatusMethodNotAllowed, Message: "method not allowed"})

	if ctx.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", ctx.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestWriteHTTPError_WithGenericError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=1.1.1.1", nil)
	ctx := ContextWithDefaults(httptest.NewRecorder(), req)

	WriteHTTPError(ctx, errors.New("boom"))

	if ctx.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", ctx.StatusCode, http.StatusInternalServerError)
	}
}

func TestWriteResponseMiddleware_WritesHeadersAndBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=1.1.1.1", nil)
	rr := httptest.NewRecorder()
	ctx := ContextWithDefaults(rr, req)
	ctx.StatusCode = http.StatusOK
	ctx.ResponseBody = map[string]string{"country": "US", "city": "New York"}
	ctx.ResponseHeader.Set("RateLimit-Limit", "10")

	WriteResponseMiddleware(ctx, func(*Context) {})

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Header().Get("RateLimit-Limit") != "10" {
		t.Fatalf("RateLimit-Limit = %q, want %q", rr.Header().Get("RateLimit-Limit"), "10")
	}

	var got map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if got["country"] != "US" || got["city"] != "New York" {
		t.Fatalf("unexpected response body: %+v", got)
	}
}
