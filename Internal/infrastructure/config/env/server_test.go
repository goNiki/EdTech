package env_test

import (
	"os"
	"testing"
	"time"

	"edtech/internal/infrastructure/config/env"
)

func TestServerConfig_CorsAllowedOrigins(t *testing.T) {
	// Set required server env vars
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_TIMEOUT", "5s")
	os.Setenv("SERVER_IDLETIMEOUT", "60s")

	t.Run("Default_Origins", func(t *testing.T) {
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		cfg, err := env.NewServerConfig()
		if err != nil {
			t.Fatalf("unexpected error creating server config: %v", err)
		}

		origins := cfg.CorsAllowedOrigins()
		if len(origins) == 0 {
			t.Fatalf("expected non-empty default origins")
		}
		expected := "http://localhost:3000"
		found := false
		for _, o := range origins {
			if o == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %q in default origins, got %v", expected, origins)
		}
	})

	t.Run("Custom_Origins", func(t *testing.T) {
		os.Setenv("CORS_ALLOWED_ORIGINS", "https://frontend.edtech.com,https://admin.edtech.com")
		defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

		cfg, err := env.NewServerConfig()
		if err != nil {
			t.Fatalf("unexpected error creating server config: %v", err)
		}

		origins := cfg.CorsAllowedOrigins()
		if len(origins) != 2 {
			t.Fatalf("expected 2 origins, got %d (%v)", len(origins), origins)
		}
		if origins[0] != "https://frontend.edtech.com" || origins[1] != "https://admin.edtech.com" {
			t.Errorf("unexpected origins: %v", origins)
		}
	})

	t.Run("Getters", func(t *testing.T) {
		cfg, err := env.NewServerConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Host() != "127.0.0.1" {
			t.Errorf("expected host 127.0.0.1, got %s", cfg.Host())
		}
		if cfg.Port() != "8080" {
			t.Errorf("expected port 8080, got %s", cfg.Port())
		}
		if cfg.Address() != "127.0.0.1:8080" {
			t.Errorf("expected address 127.0.0.1:8080, got %s", cfg.Address())
		}
		if cfg.TimeOut() != 5*time.Second {
			t.Errorf("expected timeout 5s, got %v", cfg.TimeOut())
		}
		if cfg.Idletimeout() != 60*time.Second {
			t.Errorf("expected idletimeout 60s, got %v", cfg.Idletimeout())
		}
	})
}
