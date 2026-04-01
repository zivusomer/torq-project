package findcountry

import (
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

func TestFindCountrySuccess(t *testing.T) {
	h := NewHandler(fakeStore{
		record: location.Record{Country: "US", City: "New York"},
	}, ratelimit.New(10))

	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=2.22.233.255", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestFindCountryRateLimited(t *testing.T) {
	h := NewHandler(fakeStore{
		record: location.Record{Country: "US", City: "New York"},
	}, ratelimit.New(1))

	req := httptest.NewRequest(http.MethodGet, "/v1/find-country?ip=2.22.233.255", nil)
	rr1 := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr1, req)
	h.ServeHTTP(rr2, req)

	if rr2.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rr2.Code, http.StatusTooManyRequests)
	}
}

func TestFindCountryBadRequestWhenIPMissing(t *testing.T) {
	h := NewHandler(fakeStore{}, ratelimit.New(10))

	req := httptest.NewRequest(http.MethodGet, "/v1/find-country", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}
