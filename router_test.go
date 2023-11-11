package heligo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sted/heligo"
)

type tParam struct {
	name  string
	value string
}

func TestHandle(t *testing.T) {
	tests := []struct {
		handler string
		url     string
		status  bool
		params  []tParam
	}{
		{"/a", "/a", true, nil},
		{"", "/b", false, nil},
		{"/", "/", true, nil},
		{"/app/:client_id/tokens", "/app/51/tokens", true, []tParam{{"client_id", "51"}}},
		{"/app/:client_id/tokens/:access_token", "/app/51/tokens/ky89", true, []tParam{{"client_id", "51"}, {"access_token", "ky89"}}},
		{"/bb/:gg/kk", "/bb/r/kk", true, []tParam{{"gg", "r"}}},
		{"", "/bb/rhg/kk", true, []tParam{{"gg", "rhg"}}},
		{"/a/b", "/a/b", true, nil},
		{"/aa", "/aa", true, nil},
		{"/ac", "/ac", true, nil},
		{"/boo", "/boo", true, nil},
		{"", "/aa", true, nil},
		{"", "/a/b", true, nil},
		{"", "/book", false, nil},
		{"", "/boo/u", false, nil},
		{"", "/boo", true, nil},
		{"/contribute.html", "/contribute.html", true, nil},
		{"/debugging_with_gdb.html", "/debugging_with_gdb.html", true, nil},
		{"/docs.html", "/docs.html", true, nil},
		{"/effective_go.html", "/effective_go.html", true, nil},
		{"/files.log", "/files.log", true, nil},
		{"/gccgo_contribute.html", "/gccgo_contribute.html", true, nil},
		{"/gccgo_install.html", "/gccgo_install.html", true, nil},
		{"/go-logo-black.png", "/go-logo-black.png", true, nil},
		{"/go-logo-blue.png", "/go-logo-blue.png", true, nil},
		{"/go-logo-white.png", "/go-logo-white.png", true, nil},
		{"/go1.1.html", "/go1.1.html", true, nil},
		{"/go1.2.html", "/go1.2.html", true, nil},
		{"/go1.html", "/go1.html", true, nil},
		{"/api/:test", "/api/test", true, []tParam{{"test", "test"}}},
		{"", "/api/test/other", false, nil},
		{"/api/:test/:n", "/api/test/n1", true, []tParam{{"test", "test"}, {"n", "n1"}}},
		{"/api/t/*tt", "/api/t/1", true, []tParam{{"tt", "1"}}},
		{"", "/api/t/muchmore", true, []tParam{{"tt", "muchmore"}}},
		{"", "/api/t/muchmore/andmore", true, []tParam{{"tt", "muchmore/andmore"}}},
		{"/api/*test", "/api/test/u/y", true, []tParam{{"test", "test/u/y"}}},
		{"/a/*b", "/a/b/c", true, []tParam{{"b", "b/c"}}},
		{"", "/contribute.html", true, nil},
		{"", "/debugging_with_gdb.html", true, nil},
		{"", "/docs.html", true, nil},
		{"", "/effective_go.html", true, nil},
		{"", "/files.log", true, nil},
		{"", "/gccgo_contribute.html", true, nil},
		{"", "/gccgo_install.html", true, nil},
		{"", "/go-logo-black.png", true, nil},
		{"", "/go-logo-blue.png", true, nil},
		{"", "/go-logo-white.png", true, nil},
		{"", "/go1.1.html", true, nil},
		{"", "/go1.2.html", true, nil},
		{"", "/go1.html", true, nil},
		{"", "/", true, nil},
	}

	router := heligo.New()

	for _, test := range tests {
		if test.handler != "" {
			router.Handle("GET", test.handler, func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
				if test.url != r.Request.URL.Path {
					t.Fail()
				}
				params := r.Params()
				if len(params) == len(test.params) {
					for i := 0; i < len(params); i++ {
						if params[i].Name != test.params[i].name ||
							params[i].Value != test.params[i].value {
							t.Fail()
						}
					}
				} else {
					t.Fail()
				}
				return 200, nil
			})
		}
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", test.url, nil)
		router.ServeHTTP(w, r)
		if test.status == (w.Code == http.StatusNotFound) {
			t.Fail()
		}
	}
}

func BenchmarkRouter(b *testing.B) {
	ww := httptest.NewRecorder()
	req_base, err := http.NewRequest("GET", "/base/test", nil)
	if err != nil {
		panic(err)
	}
	router := heligo.New()
	router.Handle("GET", "/base/:test", func(ctx context.Context, _ http.ResponseWriter, r heligo.Request) (int, error) {
		r.Param("test")
		return 200, nil
	})
	b.Run("Heligo base", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			router.ServeHTTP(ww, req_base)
		}
	})
}
