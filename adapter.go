package heligo

import (
	"context"
	"net/http"
)

type paramsKey struct{}

var ParamsTag = paramsKey{}

type AdapterResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *AdapterResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *AdapterResponseWriter) Status() int {
	return w.status
}

// Adapt adapts a standard http.HandlerFunc to be used as a Heligo handler.
// In the standard handler one can retrieve parameters from the context,
// using ParamsFromContext.
func Adapt(hf http.HandlerFunc) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r Request) (int, error) {
		rw := &AdapterResponseWriter{w, http.StatusOK}
		req := r.Request
		if r.params.count > 0 {
			ctx := req.Context()
			ctx = context.WithValue(ctx, ParamsTag, r.Params())
			req = req.WithContext(ctx)
		}
		hf(rw, req)
		return rw.Status(), nil
	}
}

// ParamsFromContext gets the route parameters from the context
func ParamsFromContext(ctx context.Context) []Param {
	v := ctx.Value(ParamsTag)
	if v == nil {
		return nil
	}
	return v.([]Param)
}

// AdaptMiddleware adapts a standard middleware (func(http.Handler) http.Handler)
// to be used as a Heligo middleware
func AdaptMiddleware(m func(http.Handler) http.Handler) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, w http.ResponseWriter, r Request) (int, error) {
			h := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				next(rr.Context(), ww, Request{Request: rr})
			})
			rw := &AdapterResponseWriter{w, http.StatusOK}
			req := r.Request
			m(h).ServeHTTP(rw, req)
			return rw.Status(), nil
		}
	}
}

// AdaptAsMiddleware adapts a standard http.HandlerFunc to be used as
// a Heligo middleware. Note the this is less flexible than the previous
// AdaptMiddleware, as it doesn't not allow to break the middleware chain
// except in the case of errors and calls the next handler only at the end.
func AdaptAsMiddleware(hf http.HandlerFunc) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, w http.ResponseWriter, r Request) (int, error) {
			s, err := Adapt(hf)(ctx, w, r)
			if err != nil {
				return s, err
			}
			return next(ctx, w, r)
		}
	}
}
