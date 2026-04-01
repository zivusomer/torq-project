package httpserver

import "net/http"

type Server struct {
	findCountry http.Handler
}

func New(findCountry http.Handler) *Server {
	return &Server{findCountry: findCountry}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/v1/find-country", s.findCountry)
	return mux
}
