// Package db provides database operations for the application.
package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/koban/ci/models"
)

// CreateUser creates a new user in the database.
func CreateUser(db *sql.DB, input models.CreateUserInput) (*models.User, error) {
	now := time.Now()

	result, err := db.Exec(
		"INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)",
		input.Name, input.Email, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &models.User{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetUserByID retrieves a user by ID.
func GetUserByID(db *sql.DB, id int64) (*models.User, error) {
	user := &models.User{}

	err := db.QueryRow(
		"SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database.
func GetAllUsers(db *sql.DB) ([]*models.User, error) {
	rows, err := db.Query(
		"SELECT id, name, email, created_at, updated_at FROM users ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates an existing user.
func UpdateUser(db *sql.DB, id int64, input models.UpdateUserInput) (*models.User, error) {
	user, err := GetUserByID(db, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	user.Name = input.Name
	user.Email = input.Email
	user.UpdatedAt = time.Now()

	_, err = db.Exec(
		"UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?",
		user.Name, user.Email, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user by ID.
func DeleteUser(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return nil
	}

	return nil
}
