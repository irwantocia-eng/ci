package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Defaults(t *testing.T) {
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("DB_PATH")

	cfg := LoadConfig()

	if cfg.Port != "8089" {
		t.Errorf("Expected port 8089, got %s", cfg.Port)
	}

	if cfg.DBPath != "koban.db" {
		t.Errorf("Expected db path koban.db, got %s", cfg.DBPath)
	}
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	_ = os.Setenv("PORT", "3000")
	_ = os.Setenv("DB_PATH", "test.db")
	defer func() {
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("DB_PATH")
	}()

	cfg := LoadConfig()

	if cfg.Port != "3000" {
		t.Errorf("Expected port 3000, got %s", cfg.Port)
	}

	if cfg.DBPath != "test.db" {
		t.Errorf("Expected db path test.db, got %s", cfg.DBPath)
	}
}

func TestNewServer(t *testing.T) {
	mux := http.NewServeMux()
	server := NewServer(":8080", mux)

	if server.Addr != ":8080" {
		t.Errorf("Expected addr :8080, got %s", server.Addr)
	}

	if server.ReadTimeout != 15*time.Second {
		t.Errorf("Expected ReadTimeout 15s, got %v", server.ReadTimeout)
	}

	if server.WriteTimeout != 15*time.Second {
		t.Errorf("Expected WriteTimeout 15s, got %v", server.WriteTimeout)
	}

	if server.IdleTimeout != 60*time.Second {
		t.Errorf("Expected IdleTimeout 60s, got %v", server.IdleTimeout)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler := loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestLoggingMiddleware_Headers(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	w := httptest.NewRecorder()

	var loggedMethod, loggedPath string
	handler := loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggedMethod = r.Method
		loggedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if loggedMethod != "POST" {
		t.Errorf("Expected method POST, got %s", loggedMethod)
	}
	if loggedPath != "/api/test" {
		t.Errorf("Expected path /api/test, got %s", loggedPath)
	}
}

func TestCorsMiddleware_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCorsMiddleware_OPTIONS(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()

	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin '*', got %s", allowOrigin)
	}
}

func TestCorsMiddleware_Preflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/users", nil)
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	w := httptest.NewRecorder()

	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("Expected Access-Control-Allow-Methods, got %s", methods)
	}

	headers := w.Header().Get("Access-Control-Allow-Headers")
	if headers != "Content-Type" {
		t.Errorf("Expected Access-Control-Allow-Headers Content-Type, got %s", headers)
	}
}

func TestSetupRouter(t *testing.T) {
	database, err := InitDatabase(&Config{DBPath: ":memory:"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() { _ = database.Close() }()

	mux := SetupRouter(database)

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/api/users"},
		{"POST", "/api/users"},
		{"GET", "/api/users/"},
		{"PUT", "/api/users/"},
		{"DELETE", "/api/users/"},
		{"GET", "/api/products"},
		{"POST", "/api/products"},
		{"GET", "/api/products/"},
		{"PUT", "/api/products/"},
		{"DELETE", "/api/products/"},
		{"GET", "/api/configs"},
		{"POST", "/api/configs"},
		{"GET", "/api/configs/"},
		{"PUT", "/api/configs/"},
		{"DELETE", "/api/configs/"},
		{"GET", "/"},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code == http.StatusNotFound {
			t.Errorf("Route %s %s not found", tc.method, tc.path)
		}
	}
}

func TestInitDatabase_InvalidPath(t *testing.T) {
	cfg := &Config{DBPath: "/invalid/path/that/does/not/exist/database.sqlite"}

	_, err := InitDatabase(cfg)

	if err == nil {
		t.Error("Expected error for invalid database path, got nil")
	}
}
