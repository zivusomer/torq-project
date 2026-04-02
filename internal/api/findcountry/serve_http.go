package findcountry

import (
	"net/http"

	"zivusomer/torq-project/internal/api/common"
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := common.ContextWithDefaults(w, r)
	chain := common.Chain(
		common.MiddlewareFunc(common.WriteResponseMiddleware),
		common.MiddlewareFunc(common.CallerIPIdentityMiddleware),
		common.MiddlewareFunc(common.RateLimitMiddleware),
		common.MiddlewareFunc(common.RequestIPMiddleware),
		common.MiddlewareFunc(methodGuardMiddleware),
		common.MiddlewareFunc(h.executeRequestMiddleware),
	)
	chain(ctx)
}
