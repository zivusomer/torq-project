package findcountry

import (
	"net/http"

	"zivusomer/torq-project/internal/api/common"
	"zivusomer/torq-project/internal/httpserver/httpx"
)

func methodGuardMiddleware(ctx *common.Context, next common.Next) {
	if err := httpx.RequireMethod(ctx.R, http.MethodGet); err != nil {
		common.WriteHTTPError(ctx, err)
		return
	}
	next(ctx)
}
