// Package main is the entry point for the Koban REST API server.
// Testing SonarCloud CI trigger - v5
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "koban.db"
	}

	database, err := db.NewSQLite(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

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

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprintf(w, `{"message": "Welcome to Koban API", "endpoints": {"users": "/api/users", "products": "/api/products"}}`); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})

	addr := ":" + port
	// #nosec G701,G706 - addr is from env var, logging for debug purposes
	log.Printf("Starting server on %q", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      corsMiddleware(loggingMiddleware(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
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
