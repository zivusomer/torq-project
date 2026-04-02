package findcountry

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"zivusomer/torq-project/internal/httpserver/httpx"
	"zivusomer/torq-project/internal/ratelimit"
	"zivusomer/torq-project/internal/store"
)

type Handler struct {
	store store.Resolver
}

func NewHandler(store store.Resolver) (*Handler, error) {
	if store == nil {
		return nil, fmt.Errorf("store is required")
	}

	return &Handler{
		store: store,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := httpx.RequireMethod(r, http.MethodGet); err != nil {
		writeHTTPError(w, err)
		return
	}

	if !ratelimit.Allow() {
		httpx.WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
	}

	ipParam, err := httpx.RequiredQuery(r, "ip")
	if err != nil {
		writeHTTPError(w, err)
		return
	}

	ip := net.ParseIP(ipParam)
	if ip == nil {
		httpx.WriteJSONError(w, http.StatusBadRequest, "invalid ip address")
		return
	}

	record, err := h.store.FindByIP(ip)
	if err != nil {
		if errors.Is(err, store.ErrIPNotFound) {
			httpx.WriteJSONError(w, http.StatusNotFound, "ip address not found")
			return
		}
		httpx.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, record)
}

func writeHTTPError(w http.ResponseWriter, err error) {
	var httpErr *httpx.Error
	if errors.As(err, &httpErr) {
		httpx.WriteJSONError(w, httpErr.StatusCode, httpErr.Message)
		return
	}
	httpx.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
}
