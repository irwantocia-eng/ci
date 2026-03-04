// Package handlers provides HTTP request handlers for the API endpoints.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/koban/ci/db"
	"github.com/koban/ci/models"
)

// UserListHandler handles requests to /api/users endpoint.
func UserListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllUsers(w)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// UserHandler handles requests to /api/users/{id} endpoint.
func UserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/api/users/")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getUserByID(w, id)
	case http.MethodPut:
		updateUser(w, r, id)
	case http.MethodDelete:
		deleteUser(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllUsers(w http.ResponseWriter) {
	users, err := db.GetAllUsers(dbInstance.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUserByID(w http.ResponseWriter, id int64) {
	user, err := db.GetUserByID(dbInstance.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.CreateUser(dbInstance.DB, input)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request, id int64) {
	var input models.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.UpdateUser(dbInstance.DB, id, input)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteUser(w http.ResponseWriter, id int64) {
	if err := db.DeleteUser(dbInstance.DB, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func extractID(path, prefix string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	return strconv.ParseInt(idStr, 10, 64)
}

var dbInstance *db.Database

// SetDBInstance sets the database instance for handlers to use.
func SetDBInstance(db *db.Database) {
	dbInstance = db
}
