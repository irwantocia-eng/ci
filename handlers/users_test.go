package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestUserHandler_MethodNotAllowed(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodPatch, "/api/users/1", nil)
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", res.StatusCode)
	}
}

func TestGetAllUsers_Empty(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var users []models.User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}
}

func TestCreateUser_InvalidBody(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", res.StatusCode)
	}
}

func TestUpdateUser_InvalidBody(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", res.StatusCode)
	}
}

func TestUserHandler_InvalidID(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodGet, "/api/users/abc", nil)
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", res.StatusCode)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := db.CreateUser(database.DB, models.CreateUserInput{
		Name:  "First User",
		Email: "duplicate@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	userInput := models.CreateUserInput{
		Name:  "Second User",
		Email: "duplicate@example.com",
	}
	body, _ := json.Marshal(userInput)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", res.StatusCode)
	}
}

func TestUpdateUser_DuplicateEmail(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := db.CreateUser(database.DB, models.CreateUserInput{
		Name:  "First User",
		Email: "existing@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	createdUser, err := db.CreateUser(database.DB, models.CreateUserInput{
		Name:  "Second User",
		Email: "second@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create second user: %v", err)
	}

	updateInput := models.UpdateUserInput{
		Name:  "Updated Name",
		Email: "existing@example.com",
	}
	body, _ := json.Marshal(updateInput)

	req := httptest.NewRequest(http.MethodPut, "/api/users/"+fmt.Sprintf("%d", createdUser.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", res.StatusCode)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	updateInput := models.UpdateUserInput{
		Name:  "Test",
		Email: "test@example.com",
	}
	body, _ := json.Marshal(updateInput)

	req := httptest.NewRequest(http.MethodPut, "/api/users/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", res.StatusCode)
	}
}

func TestCreateUser_EmptyName(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	userInput := models.CreateUserInput{
		Name:  "",
		Email: "valid@example.com",
	}
	body, _ := json.Marshal(userInput)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", res.StatusCode)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodDelete, "/api/users/999", nil)
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", res.StatusCode)
	}
}

func TestGetUserByID_DBError(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdUser, err := db.CreateUser(database.DB, models.CreateUserInput{
		Name:  "Test User",
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	_ = database.Close()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%d", createdUser.ID), nil)
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", res.StatusCode)
	}
}

func TestGetAllUsers_DBError(t *testing.T) {
	database := setupTestDB(t)
	_ = database.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", res.StatusCode)
	}
}

func TestCreateUser_DBError(t *testing.T) {
	database := setupTestDB(t)
	_ = database.Close()

	userInput := models.CreateUserInput{
		Name:  "Test User",
		Email: "test@example.com",
	}
	body, _ := json.Marshal(userInput)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserListHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", res.StatusCode)
	}
}

func TestUpdateUser_DBError(t *testing.T) {
	database := setupTestDB(t)
	_ = database.Close()

	updateInput := models.UpdateUserInput{
		Name:  "Test",
		Email: "test@example.com",
	}
	body, _ := json.Marshal(updateInput)

	req := httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", res.StatusCode)
	}
}

func TestDeleteUser_DBError(t *testing.T) {
	database := setupTestDB(t)
	_ = database.Close()

	req := httptest.NewRequest(http.MethodDelete, "/api/users/1", nil)
	w := httptest.NewRecorder()

	UserHandler(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", res.StatusCode)
	}
}
