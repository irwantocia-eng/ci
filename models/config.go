package models

import "time"

// Config represents a key-value configuration stored in the database.
type Config struct {
	ID        int64     `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateConfigInput represents the input for creating a new config.
type CreateConfigInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// UpdateConfigInput represents the input for updating an existing config.
type UpdateConfigInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
