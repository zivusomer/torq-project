package common

import (
	"net"
	"net/http"

	"zivusomer/torq-project/internal/httpserver/httpx"
)

func RequestIPMiddleware(ctx *Context, next Next) {
	ipParam, err := httpx.RequiredQuery(ctx.R, "ip")
	if err != nil {
		WriteHTTPError(ctx, err)
		return
	}

	ip, err := parseRequestedIP(ipParam)
	if err != nil {
		WriteHTTPError(ctx, err)
		return
	}

	ctx.RequestedIP = ip
	next(ctx)
}

func parseRequestedIP(ipParam string) (net.IP, error) {
	ip := net.ParseIP(ipParam)
	if ip == nil {
		return nil, &httpx.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid ip address",
		}
	}

	return ip, nil
}
