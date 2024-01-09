package heligo_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sted/heligo"
)

func TestAdapter(t *testing.T) {
	router := heligo.New()
	standardHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("Test adapter"))
	}
	router.Handle("GET", "/adapt", heligo.AdaptFunc(standardHandler))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/adapt", nil)
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusTeapot {
		t.Fail()
	}
	content, _ := io.ReadAll(w.Result().Body)
	w.Result().Body.Close()
	if !bytes.Equal(content, []byte("Test adapter")) {
		t.Fail()
	}
}

func TestAdapterParams(t *testing.T) {
	router := heligo.New()
	standardHandler := func(w http.ResponseWriter, r *http.Request) {
		params := heligo.ParamsFromContext(r.Context())
		if r.URL.Path == "/adapt/none" {
			if len(params) != 0 {
				t.Fail()
			}
		} else if len(params) != 3 ||
			params[0].Name != "one" || params[0].Value != "1" ||
			params[1].Name != "two" || params[1].Value != "2" ||
			params[2].Name != "three" || params[2].Value != "3/45" {
			t.Fail()
		}
	}
	router.Handle("GET", "/adapt/:one/:two/*three", heligo.AdaptFunc(standardHandler))
	router.Handle("GET", "/adapt/none", heligo.AdaptFunc(standardHandler))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/adapt/1/2/3/45", nil)
	router.ServeHTTP(w, r)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/adapt/none", nil)
	router.ServeHTTP(w, r)
}

func TestMiddlewareAdapter(t *testing.T) {
	router := heligo.New()
	router.Use(heligo.AdaptMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/adapt/m1" {
				w.WriteHeader(http.StatusTeapot)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}))
	h := func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
		return 200, nil
	}
	router.Handle("GET", "/adapt/*m", h)
	w := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/adapt/m1", nil)
	router.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusTeapot {
		t.Fail()
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/adapt/m2", nil)
	router.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Fail()
	}
}

func TestMiddlewareFuncAdapter(t *testing.T) {
	router := heligo.New()
	router.Use(heligo.AdaptFuncAsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		if r.URL.Path == "/adapt/m1" {
			headers["Access-Control-Allow-Origin"] = []string{"*"}
		} else {
			headers["Access-Control-Allow-Origin"] = r.Header["Origin"]
		}
	}))
	h := func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
		return 200, nil
	}
	router.Handle("GET", "/adapt/*m", h)
	w := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/adapt/m1", nil)
	router.ServeHTTP(w, r)
	if len(w.Result().Header["Access-Control-Allow-Origin"]) != 1 {
		t.Fail()
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/adapt/m2", nil)
	router.ServeHTTP(w, r)
	if len(w.Result().Header["Access-Control-Allow-Origin"]) != 0 {
		t.Fail()
	}
}
