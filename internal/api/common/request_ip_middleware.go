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

	ip := net.ParseIP(ipParam)
	if ip == nil {
		ctx.StatusCode = http.StatusBadRequest
		ctx.ResponseBody = map[string]string{"error": "invalid ip address"}
		return
	}

	ctx.RequestedIP = ip
	next(ctx)
}
