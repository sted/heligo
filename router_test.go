package heligo_test

import (
	"context"
	"io"
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

	var gotParams []heligo.Param
	for _, test := range tests {
		if test.handler != "" {
			router.Handle("GET", test.handler, func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
				gotParams = r.Params()
				return 200, nil
			})
		}
		gotParams = nil
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", test.url, nil)
		router.ServeHTTP(w, r)
		if test.status == (w.Code == http.StatusNotFound) {
			t.Fatalf("url %s: expected found=%v, got status %d", test.url, test.status, w.Code)
		}
		if !test.status {
			continue
		}
		if len(gotParams) != len(test.params) {
			t.Fatalf("url %s: expected %d params, got %d", test.url, len(test.params), len(gotParams))
		}
		for i, p := range gotParams {
			if p.Name != test.params[i].name || p.Value != test.params[i].value {
				t.Fatalf("url %s: param %d: expected %s=%s, got %s=%s",
					test.url, i, test.params[i].name, test.params[i].value, p.Name, p.Value)
			}
		}
	}
}

// helper straight from go-chi/chi
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	return resp, string(respBody)
}

func TestHead(t *testing.T) {
	router := heligo.New()
	router.Handle("GET", "/head", func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
		w.Header().Set("X-Test", "yes")
		w.Write([]byte("test"))
		return 200, nil
	})
	ts := httptest.NewServer(router)
	defer ts.Close()
	resp, body := testRequest(t, ts, "HEAD", "/head", nil)
	if resp.StatusCode == http.StatusNotFound || string(body) != "" || resp.Header.Get("X-Test") != "yes" {
		t.Fail()
	}
	resp, body = testRequest(t, ts, "GET", "/head", nil)
	if resp.StatusCode == http.StatusNotFound || string(body) != "test" || resp.Header.Get("X-Test") != "yes" {
		t.Fail()
	}
}

func TestTrailingSlash(t *testing.T) {
	router := heligo.New()
	router.TrailingSlash = true

	handler := func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
		return 200, nil
	}

	// Register without trailing slash
	router.Handle("GET", "/users", handler)
	// Register with trailing slash
	router.Handle("GET", "/posts/", handler)
	// Root should not be duplicated
	router.Handle("GET", "/", handler)
	// Param route
	router.Handle("GET", "/users/:id", handler)
	// Wildcard should NOT get trailing slash
	router.Handle("GET", "/static/*filepath", handler)
	// Param + static suffix
	router.Handle("GET", "/users/:id/profile", handler)

	tests := []struct {
		url    string
		status int
	}{
		{"/users", 200},
		{"/users/", 200},
		{"/posts", 200},
		{"/posts/", 200},
		{"/", 200},
		{"/users/42", 200},
		{"/users/42/", 200},
		{"/static/js/app.js", 200},
		{"/users/42/profile", 200},
		{"/users/42/profile/", 200},
		{"/nonexistent", 404},
		{"/nonexistent/", 404},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", test.url, nil)
		router.ServeHTTP(w, r)
		if w.Code != test.status {
			t.Errorf("url %s: expected %d, got %d", test.url, test.status, w.Code)
		}
	}
}

func TestTrailingSlashOff(t *testing.T) {
	router := heligo.New()
	// TrailingSlash defaults to false

	handler := func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
		return 200, nil
	}

	router.Handle("GET", "/users", handler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/", nil)
	router.ServeHTTP(w, r)
	if w.Code != 404 {
		t.Errorf("expected 404 for /users/ with TrailingSlash off, got %d", w.Code)
	}
}

func BenchmarkRouter(b *testing.B) {
	ww := httptest.NewRecorder()
	req_base, err := http.NewRequest("GET", "/base/test", nil)
	if err != nil {
		b.Fatal(err)
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
