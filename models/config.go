package models

import "time"

type Config struct {
	ID        int64     `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateConfigInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateConfigInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
