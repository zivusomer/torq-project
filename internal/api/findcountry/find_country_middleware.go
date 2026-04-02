package findcountry

import (
	"errors"
	"net/http"

	"zivusomer/torq-project/internal/api/common"
	"zivusomer/torq-project/internal/store"
)

func (h *Handler) FindCountryMiddleware(ctx *common.Context, next common.Next) {
	record, err := h.store.FindByIP(ctx.RequestedIP)
	if err != nil {
		if errors.Is(err, store.ErrIPNotFound) {
			ctx.StatusCode = http.StatusNotFound
			ctx.ResponseBody = map[string]string{"error": "ip address not found"}
			return
		}
		ctx.StatusCode = http.StatusInternalServerError
		ctx.ResponseBody = map[string]string{"error": "internal server error"}
		return
	}

	ctx.StatusCode = http.StatusOK
	ctx.ResponseBody = record
	next(ctx)
}
