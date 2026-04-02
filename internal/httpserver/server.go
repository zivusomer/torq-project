package httpserver

import (
	"fmt"
	"net/http"
)

type Route struct {
	Path    string
	Handler http.Handler
}

type Server struct {
	routes []Route
}

func New(routes []Route) (*Server, error) {
	if len(routes) == 0 {
		return nil, fmt.Errorf("at least one route is required")
	}

	for _, route := range routes {
		if route.Path == "" {
			return nil, fmt.Errorf("route path is required")
		}
		if route.Handler == nil {
			return nil, fmt.Errorf("route handler is required for path %s", route.Path)
		}
	}

	return &Server{routes: routes}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	for _, route := range s.routes {
		mux.Handle(route.Path, route.Handler)
	}
	return mux
}
