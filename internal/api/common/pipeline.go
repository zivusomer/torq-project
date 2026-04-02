package common

import (
	"net"
	"net/http"
)

type Context struct {
	W              http.ResponseWriter
	R              *http.Request
	CallerKey      string
	RequestedIP    net.IP
	StatusCode     int
	ResponseBody   any
	ResponseHeader http.Header
}

type Next func(*Context)

type Middleware interface {
	Handle(*Context, Next)
}

type MiddlewareFunc func(*Context, Next)

func (f MiddlewareFunc) Handle(ctx *Context, next Next) {
	f(ctx, next)
}

func ContextWithDefaults(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:              w,
		R:              r,
		StatusCode:     http.StatusInternalServerError,
		ResponseBody:   map[string]string{"error": "internal server error"},
		ResponseHeader: make(http.Header),
	}
}

func Chain(middlewares ...Middleware) Next {
	final := func(*Context) {}
	for i := len(middlewares) - 1; i >= 0; i-- {
		current := middlewares[i]
		next := final
		final = func(ctx *Context) {
			current.Handle(ctx, next)
		}
	}
	return final
}
