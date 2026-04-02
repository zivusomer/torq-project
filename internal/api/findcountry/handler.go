package findcountry

import (
	"fmt"

	"zivusomer/torq-project/internal/store"
)

func NewHandler(store store.Resolver) (*Handler, error) {
	if store == nil {
		return nil, fmt.Errorf("store is required")
	}

	return &Handler{
		store: store,
	}, nil
}
