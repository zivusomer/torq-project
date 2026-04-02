package findcountry

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"zivusomer/torq-project/internal/location"
	"zivusomer/torq-project/internal/ratelimit"
)

type fakeStore struct {
	record location.Record
	err    error
}

func (f fakeStore) FindByIP(_ net.IP) (location.Record, error) {
	if f.err != nil {
		return location.Record{}, f.err
	}
	return f.record, nil
}

func TestFindCountryValidIPAddress(t *testing.T) {
	h := setupHandler(t, 10, fakeStore{
		record: location.Record{Country: "US", City: "New York"},
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=2.22.233.255", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var got location.Record
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if got.Country != "US" || got.City != "New York" {
		t.Fatalf("response body = %+v, want country=US city=New York", got)
	}
}

func TestFindCountryInvalidIPAddress(t *testing.T) {
	h := setupHandler(t, 10, fakeStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=not-an-ip", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestFindCountryRateLimited(t *testing.T) {
	h := setupHandler(t, 1, fakeStore{
		record: location.Record{Country: "US", City: "New York"},
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=2.22.233.255", nil)
	rr1 := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr1, req)
	h.ServeHTTP(rr2, req)

	if rr2.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rr2.Code, http.StatusTooManyRequests)
	}
	if rr2.Header().Get("Retry-After") == "" {
		t.Fatalf("expected Retry-After header when rate limited")
	}
	if rr2.Header().Get("RateLimit-Limit") == "" {
		t.Fatalf("expected RateLimit-Limit header")
	}
	if rr2.Header().Get("RateLimit-Remaining") == "" {
		t.Fatalf("expected RateLimit-Remaining header")
	}
	if rr2.Header().Get("RateLimit-Reset") == "" {
		t.Fatalf("expected RateLimit-Reset header")
	}
}

func TestFindCountryBadRequestWhenIPMissing(t *testing.T) {
	h := setupHandler(t, 10, fakeStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/find-country", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestFindCountryHandlerRejectsNilDeps(t *testing.T) {
	_, err := NewHandler(nil)
	if err == nil {
		t.Fatalf("expected error for nil dependencies")
	}
}

func setupHandler(t *testing.T, rateLimit int, store fakeStore) *Handler {
	t.Helper()

	if err := ratelimit.Init(rateLimit); err != nil {
		t.Fatalf("ratelimit.Init() error: %v", err)
	}

	h, err := NewHandler(store)
	if err != nil {
		t.Fatalf("NewHandler() error: %v", err)
	}

	return h
}
