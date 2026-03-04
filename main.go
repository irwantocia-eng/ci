// Package main is the entry point for the Koban REST API server.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/koban/ci/db"
	"github.com/koban/ci/handlers"
)

// Config holds server configuration
type Config struct {
	Port   string
	DBPath string
}

// LoadConfig reads configuration from environment variables
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "koban.db"
	}

	return &Config{
		Port:   port,
		DBPath: dbPath,
	}
}

// InitDatabase initializes the database and runs migrations
func InitDatabase(cfg *Config) (*db.Database, error) {
	database, err := db.NewSQLite(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := database.Migrate(); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}

// SetupRouter creates and configures the HTTP router
func SetupRouter(database *db.Database) *http.ServeMux {
	handlers.SetDBInstance(database)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/users", handlers.UserListHandler)
	mux.HandleFunc("POST /api/users", handlers.UserListHandler)
	mux.HandleFunc("GET /api/users/", handlers.UserHandler)
	mux.HandleFunc("PUT /api/users/", handlers.UserHandler)
	mux.HandleFunc("DELETE /api/users/", handlers.UserHandler)

	mux.HandleFunc("GET /api/products", handlers.ProductListHandler)
	mux.HandleFunc("POST /api/products", handlers.ProductListHandler)
	mux.HandleFunc("GET /api/products/", handlers.ProductHandler)
	mux.HandleFunc("PUT /api/products/", handlers.ProductHandler)
	mux.HandleFunc("DELETE /api/products/", handlers.ProductHandler)

	mux.HandleFunc("GET /api/configs", handlers.ConfigListHandler)
	mux.HandleFunc("POST /api/configs", handlers.ConfigListHandler)
	mux.HandleFunc("GET /api/configs/", handlers.ConfigHandler)
	mux.HandleFunc("PUT /api/configs/", handlers.ConfigHandler)
	mux.HandleFunc("DELETE /api/configs/", handlers.ConfigHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprintf(w, `{"message": "Welcome to Koban API", "endpoints": {"users": "/api/users", "products": "/api/products", "configs": "/api/configs"}}`); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})

	return mux
}

// NewServer creates a new HTTP server with the given configuration
func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      corsMiddleware(loggingMiddleware(handler)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// RunServer starts the HTTP server (typically not called in tests)
func RunServer(srv *http.Server) error {
	log.Printf("Starting server on %q", srv.Addr)
	return srv.ListenAndServe()
}

func main() {
	cfg := LoadConfig()

	database, err := InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	mux := SetupRouter(database)

	addr := ":" + cfg.Port
	server := NewServer(addr, mux)

	if err := RunServer(server); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// #nosec G706 - logging request info for debugging
		log.Printf("%s %q", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
