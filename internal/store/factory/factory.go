package factory

import (
	"fmt"

	"zivusomer/torq-project/internal/store"
	"zivusomer/torq-project/internal/store/csvstore"
)

func New(datastoreType, path string) (store.Resolver, error) {
	switch datastoreType {
	case "csv":
		return csvstore.New(path)
	default:
		return nil, fmt.Errorf("unsupported datastore type: %s", datastoreType)
	}
}
