package store

import (
	"errors"
	"net"

	"zivusomer/torq-project/internal/location"
)

var ErrIPNotFound = errors.New("ip not found")

type Resolver interface {
	FindByIP(ip net.IP) (location.Record, error)
}
