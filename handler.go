package heligo

import (
	"context"
	"fmt"
	"net/http"
	"path"
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

// Recover returns a middleware that recovers from panics in downstream handlers.
// If a panic occurs, it responds with 500 Internal Server Error.
// The optional onPanic callback receives the recovered value.
func Recover(onPanic func(v any)) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, w http.ResponseWriter, r Request) (status int, err error) {
			defer func() {
				if v := recover(); v != nil {
					if onPanic != nil {
						onPanic(v)
					}
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					status = http.StatusInternalServerError
					err = fmt.Errorf("panic: %v", v)
				}
			}()
			return next(ctx, w, r)
		}
	}
}

// CleanPaths returns a middleware that cleans URL paths
// containing //, /./ or /../ sequences using path.Clean.
func CleanPaths() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, w http.ResponseWriter, r Request) (int, error) {
			if needsClean(r.URL.Path) {
				r.URL.Path = path.Clean(r.URL.Path)
			}
			return next(ctx, w, r)
		}
	}
}
