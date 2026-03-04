package db

import (
	"testing"

	"github.com/koban/ci/models"
)

func TestNewSQLite_InvalidPath(t *testing.T) {
	_, err := NewSQLite("/invalid/path/that/does/not/exist/db.sqlite")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestMigrate_InvalidDB(t *testing.T) {
	db, err := NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to create db: %v", err)
	}
	_ = db.Close()

	err = db.Migrate()
	if err == nil {
		t.Error("Expected error when migrating closed db")
	}
}

func setupTestDB(t *testing.T) *Database {
	database, err := NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return database
}

func TestCreateUser(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	user, err := CreateUser(database.DB, models.CreateUserInput{
		Name:  "Test User",
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}

	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.Name)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestGetUserByID(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdUser, err := CreateUser(database.DB, models.CreateUserInput{
		Name:  "Get Test User",
		Email: "get@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err := GetUserByID(database.DB, createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user == nil {
		t.Fatal("Expected user to be found")
	}

	if user.ID != createdUser.ID {
		t.Errorf("Expected ID %d, got %d", createdUser.ID, user.ID)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	user, err := GetUserByID(database.DB, 999)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user != nil {
		t.Error("Expected nil user for non-existent ID")
	}
}

func TestGetAllUsers(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := CreateUser(database.DB, models.CreateUserInput{
		Name:  "User 1",
		Email: "user1@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}

	_, err = CreateUser(database.DB, models.CreateUserInput{
		Name:  "User 2",
		Email: "user2@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	users, err := GetAllUsers(database.DB)
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestUpdateUser(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdUser, err := CreateUser(database.DB, models.CreateUserInput{
		Name:  "Original Name",
		Email: "original@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	updatedUser, err := UpdateUser(database.DB, createdUser.ID, models.UpdateUserInput{
		Name:  "Updated Name",
		Email: "updated@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if updatedUser.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", updatedUser.Name)
	}

	if updatedUser.Email != "updated@example.com" {
		t.Errorf("Expected email 'updated@example.com', got '%s'", updatedUser.Email)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	user, err := UpdateUser(database.DB, 999, models.UpdateUserInput{
		Name:  "Test",
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if user != nil {
		t.Error("Expected nil user for non-existent ID")
	}
}

func TestDeleteUser(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdUser, err := CreateUser(database.DB, models.CreateUserInput{
		Name:  "Delete Me",
		Email: "delete@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	err = DeleteUser(database.DB, createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	user, _ := GetUserByID(database.DB, createdUser.ID)
	if user != nil {
		t.Error("Expected user to be deleted")
	}
}
