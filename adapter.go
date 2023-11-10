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

func ParamsFromContext(ctx context.Context) []Param {
	return ctx.Value(ParamsTag).([]Param)
}
