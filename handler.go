package heligo

import (
	"context"
	"net/http"
)

// Handler is the signature of a Heligo handler.
// It gets a standard context taken from http.Request, a http.RequestWriter and
// a heligo.Request.
// You should return the HTTP status code and a potential error (or nil).
type Handler func(ctx context.Context, w http.ResponseWriter, r Request) (int, error)

// Middleware is the signature of a Heligo middleware.
type Middleware func(Handler) Handler

func chain(h Handler, middlewares []Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
