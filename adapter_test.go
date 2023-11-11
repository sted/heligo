package heligo_test

import (
	"bytes"
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
	router.Handle("GET", "/adapt", heligo.Adapt(standardHandler))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/adapt", nil)
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusTeapot {
		t.Fail()
	}
	body, _ := io.ReadAll(w.Result().Body)
	w.Result().Body.Close()
	if !bytes.Equal(body, []byte("Test adapter")) {
		t.Fail()
	}
}

func TestAdapterParams(t *testing.T) {
	router := heligo.New()
	standardHandler := func(w http.ResponseWriter, r *http.Request) {
		params := heligo.ParamsFromContext(r.Context())
		if len(params) != 3 ||
			params[0].Name != "one" || params[0].Value != "1" ||
			params[1].Name != "two" || params[1].Value != "2" ||
			params[2].Name != "three" || params[2].Value != "3/45" {
			t.Fail()
		}
	}
	router.Handle("GET", "/adapt/:one/:two/*three", heligo.Adapt(standardHandler))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/adapt/1/2/3/45", nil)
	router.ServeHTTP(w, r)
}
