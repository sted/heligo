package heligo_test

import (
	"bytes"
	"context"
	"heligo"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		url  string
		body []byte
	}{
		{"/mw", []byte("m0 mw ")},
		{"/g1/mw", []byte("m0 m1 mw e1 ")},
		{"/g2/mw", []byte("m0 m2 m1 mw e1 e2 ")},
		{"/g1/g3/mw", []byte("m0 m1 m2 mw e2 e1 ")},
		{"/g2/g4/mw", []byte("m0 m2 m1 m1 m2 mw e2 e1 e1 e2 ")},
	}

	router := heligo.New()
	m0 := func(next heligo.Handler) heligo.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
			w.Write([]byte("m0 "))
			return next(ctx, w, r)
		}
	}
	router.Use(m0)
	m1 := func(next heligo.Handler) heligo.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
			w.Write([]byte("m1 "))
			status, err := next(ctx, w, r)
			w.Write([]byte("e1 "))
			return status, err
		}
	}
	m2 := func(next heligo.Handler) heligo.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
			w.Write([]byte("m2 "))
			status, err := next(ctx, w, r)
			w.Write([]byte("e2 "))
			return status, err
		}
	}
	h := func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
		w.Write([]byte("mw "))
		return 0, nil
	}
	g1 := router.Group("/g1", m1)
	g2 := router.Group("/g2", m2, m1)
	g3 := g1.Group("/g3", m2)
	g4 := g2.Group("/g4", m1, m2)
	router.Handle("GET", "/mw", h)
	g1.Handle("GET", "/mw", h)
	g2.Handle("GET", "/mw", h)
	g3.Handle("GET", "/mw", h)
	g4.Handle("GET", "/mw", h)

	for _, test := range tests {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", test.url, nil)
		router.ServeHTTP(w, r)
		body, _ := io.ReadAll(w.Result().Body)
		w.Result().Body.Close()
		if !bytes.Equal(body, test.body) {
			t.Fail()
		}
	}
}
