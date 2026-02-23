package heligo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sted/heligo"
)

// noop handler used by all benchmarks - minimal work to isolate routing cost
func benchHandler(_ context.Context, _ http.ResponseWriter, _ heligo.Request) (int, error) {
	return 200, nil
}

func benchHandlerParam(_ context.Context, _ http.ResponseWriter, r heligo.Request) (int, error) {
	r.Param("id")
	return 200, nil
}

// buildLargeRouter creates a router with many routes to simulate a real application
func buildLargeRouter() *heligo.Router {
	router := heligo.New()

	// Static routes (typical API surface)
	statics := []string{
		"/",
		"/favicon.ico",
		"/health",
		"/api/v1/status",
		"/api/v1/users",
		"/api/v1/users/search",
		"/api/v1/projects",
		"/api/v1/projects/search",
		"/api/v1/organizations",
		"/api/v1/organizations/settings",
		"/api/v1/billing",
		"/api/v1/billing/invoices",
		"/api/v1/notifications",
		"/api/v1/notifications/preferences",
		"/api/v2/users",
		"/api/v2/projects",
		"/docs",
		"/docs/api",
		"/docs/guides",
		"/docs/tutorials",
		"/about",
		"/contact",
		"/login",
		"/logout",
		"/signup",
		"/settings",
		"/settings/profile",
		"/settings/security",
		"/settings/tokens",
	}
	for _, p := range statics {
		router.Handle("GET", p, benchHandler)
	}

	// Parameterized routes
	router.Handle("GET", "/api/v1/users/:id", benchHandlerParam)
	router.Handle("GET", "/api/v1/users/:id/projects", benchHandlerParam)
	router.Handle("GET", "/api/v1/users/:id/tokens", benchHandlerParam)
	router.Handle("GET", "/api/v1/projects/:id", benchHandlerParam)
	router.Handle("GET", "/api/v1/projects/:id/members", benchHandlerParam)
	router.Handle("GET", "/api/v1/projects/:id/settings", benchHandlerParam)
	router.Handle("GET", "/api/v1/organizations/:id", benchHandlerParam)
	router.Handle("GET", "/api/v1/organizations/:id/members", benchHandlerParam)
	router.Handle("POST", "/api/v1/users", benchHandler)
	router.Handle("POST", "/api/v1/projects", benchHandler)
	router.Handle("PUT", "/api/v1/users/:id", benchHandlerParam)
	router.Handle("PUT", "/api/v1/projects/:id", benchHandlerParam)
	router.Handle("DELETE", "/api/v1/users/:id", benchHandlerParam)
	router.Handle("DELETE", "/api/v1/projects/:id", benchHandlerParam)

	// Multi-param routes
	router.Handle("GET", "/api/v1/projects/:pid/members/:uid", benchHandler)
	router.Handle("GET", "/api/v1/organizations/:oid/projects/:pid", benchHandler)
	router.Handle("GET", "/api/v1/organizations/:oid/projects/:pid/members/:uid", benchHandler)

	// Wildcard routes
	router.Handle("GET", "/static/*filepath", benchHandler)
	router.Handle("GET", "/files/:user/*filepath", benchHandler)

	return router
}

var w = httptest.NewRecorder()

// --- Static route benchmarks ---

func BenchmarkStaticRoot(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkStaticShort(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/health", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkStaticDeep(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/api/v1/notifications/preferences", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

// --- Param route benchmarks ---

func BenchmarkParam1(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/api/v1/users/42", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkParam1Long(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/api/v1/users/42/projects", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkParam2(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/api/v1/projects/99/members/42", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkParam3(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/api/v1/organizations/5/projects/99/members/42", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

// --- Wildcard benchmarks ---

func BenchmarkWildcard(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/static/js/app.min.js", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkWildcardParamMixed(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/files/john/documents/report.pdf", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

// --- Other methods ---

func BenchmarkPOST(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("POST", "/api/v1/users", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkPUT(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("PUT", "/api/v1/projects/99", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

// --- Miss benchmark ---

func BenchmarkNotFound(b *testing.B) {
	router := buildLargeRouter()
	req, _ := http.NewRequest("GET", "/api/v1/nonexistent/path", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

// --- Middleware overhead ---

func BenchmarkMiddleware3(b *testing.B) {
	router := heligo.New()
	noop := func(next heligo.Handler) heligo.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
			return next(ctx, w, r)
		}
	}
	router.Use(noop, noop, noop)
	router.Handle("GET", "/api/v1/users/:id", benchHandlerParam)
	req, _ := http.NewRequest("GET", "/api/v1/users/42", nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}
