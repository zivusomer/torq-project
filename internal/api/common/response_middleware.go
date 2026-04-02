package common

import "zivusomer/torq-project/internal/httpserver/httpx"

func WriteResponseMiddleware(ctx *Context, next Next) {
	next(ctx)
	for key, values := range ctx.ResponseHeader {
		for _, value := range values {
			ctx.W.Header().Add(key, value)
		}
	}
	httpx.WriteJSON(ctx.W, ctx.StatusCode, ctx.ResponseBody)
}
