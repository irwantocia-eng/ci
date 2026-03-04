package models

import "time"

// User represents a user in the system.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserInput represents the input for creating a new user.
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserInput represents the input for updating an existing user.
type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
