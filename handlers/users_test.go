package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koban/ci/db"
	"github.com/koban/ci/models"
)

func setupTestDB(t *testing.T) *db.Database {
	database, err := db.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	SetDBInstance(database)
	return database
}

func TestUserListHandler_GET(t *testing.T) {
	database := setupTestDB(t)
	defer func() {
		_ = database.Close()
	}()

	_, err := db.CreateUser(database.DB, models.CreateUserInput{
		Name:  "Test User",
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var users []models.User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	if users[0].Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", users[0].Name)
	}
}

func TestUserListHandler_POST(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	userInput := models.CreateUserInput{
		Name:  "New User",
		Email: "new@example.com",
	}
	body, _ := json.Marshal(userInput)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", res.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if user.Name != userInput.Name {
		t.Errorf("Expected name '%s', got '%s'", userInput.Name, user.Name)
	}

	if user.Email != userInput.Email {
		t.Errorf("Expected email '%s', got '%s'", userInput.Email, user.Email)
	}
}

func TestUserHandler_GET(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdUser, err := db.CreateUser(database.DB, models.CreateUserInput{
		Name:  "Get User",
		Email: "get@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if user.ID != createdUser.ID {
		t.Errorf("Expected user ID %d, got %d", createdUser.ID, user.ID)
	}
}

func TestUserHandler_GET_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodGet, "/api/users/999", nil)
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", res.StatusCode)
	}
}

func TestUserHandler_PUT(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := db.CreateUser(database.DB, models.CreateUserInput{
		Name:  "Original",
		Email: "original@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	updateInput := models.UpdateUserInput{
		Name:  "Updated",
		Email: "updated@example.com",
	}
	body, _ := json.Marshal(updateInput)

	req := httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if user.Name != updateInput.Name {
		t.Errorf("Expected name '%s', got '%s'", updateInput.Name, user.Name)
	}
}

func TestUserHandler_DELETE(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := db.CreateUser(database.DB, models.CreateUserInput{
		Name:  "Delete User",
		Email: "delete@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/users/1", nil)
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", res.StatusCode)
	}

	// Verify user was deleted
	user, err := db.GetUserByID(database.DB, 1)
	if err != nil {
		t.Fatalf("Failed to check user: %v", err)
	}

	if user != nil {
		t.Error("Expected user to be deleted, but user still exists")
	}
}

func TestUserListHandler_MethodNotAllowed(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodPatch, "/api/users", nil)
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", res.StatusCode)
	}
}
