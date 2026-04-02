package common

import (
	"errors"
	"net/http"

	"zivusomer/torq-project/internal/httpserver/httpx"
)

func WriteHTTPError(ctx *Context, err error) {
	var httpErr *httpx.Error
	if errors.As(err, &httpErr) {
		ctx.StatusCode = httpErr.StatusCode
		ctx.ResponseBody = map[string]string{"error": httpErr.Message}
		return
	}
	ctx.StatusCode = http.StatusInternalServerError
	ctx.ResponseBody = map[string]string{"error": "internal server error"}
}
