package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func TestCORS_PreflightAndValidation(t *testing.T) {
	r := chi.NewRouter()
	allowedOrigins := []string{"http://localhost:3000", "https://edtech.example.com"}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	t.Run("Allowed Origin Preflight OPTIONS", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Authorization, Content-Type")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
			t.Errorf("expected Access-Control-Allow-Origin 'http://localhost:3000', got %q", got)
		}
		if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
			t.Errorf("expected Access-Control-Allow-Credentials 'true', got %q", got)
		}
	})

	t.Run("Allowed Origin Actual Request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Origin", "https://edtech.example.com")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://edtech.example.com" {
			t.Errorf("expected Access-Control-Allow-Origin 'https://edtech.example.com', got %q", got)
		}
		if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
			t.Errorf("expected Access-Control-Allow-Credentials 'true', got %q", got)
		}
	})

	t.Run("Disallowed Origin Rejected", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Origin", "http://evil-attacker.com")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("expected empty Access-Control-Allow-Origin for unauthorized origin, got %q", got)
		}
	})
}
