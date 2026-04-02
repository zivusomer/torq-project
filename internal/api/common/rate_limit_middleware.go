package common

import (
	"net/http"
	"strconv"

	"zivusomer/torq-project/internal/logging"
	"zivusomer/torq-project/internal/ratelimit"
)

func RateLimitMiddleware(ctx *Context, next Next) {
	key := ctx.CallerKey
	if key == "" {
		logging.Logger.Error("missing caller key in rate limiter context, falling back to anonymous")
		key = "anonymous"
	}
	decision := ratelimit.AllowForKey(key)
	writeRateLimitHeaders(ctx.ResponseHeader, decision)
	if !decision.Allowed {
		ctx.ResponseHeader.Set("Retry-After", strconv.Itoa(decision.RetryAfterSeconds))
		ctx.StatusCode = http.StatusTooManyRequests
		ctx.ResponseBody = map[string]string{"error": "rate limit exceeded"}
		return
	}

	next(ctx)
}

func writeRateLimitHeaders(headers http.Header, decision ratelimit.Decision) {
	headers.Set("RateLimit-Limit", strconv.Itoa(decision.Limit))
	headers.Set("RateLimit-Remaining", strconv.Itoa(decision.Remaining))
	headers.Set("RateLimit-Reset", strconv.Itoa(decision.ResetSeconds))
}
