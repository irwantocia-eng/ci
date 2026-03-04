package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/koban/ci/models"
)

func CreateConfig(db *sql.DB, input models.CreateConfigInput) (*models.Config, error) {
	now := time.Now()

	result, err := db.Exec(
		"INSERT INTO configs (key, value, created_at, updated_at) VALUES (?, ?, ?, ?)",
		input.Key, input.Value, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &models.Config{
		ID:        id,
		Key:       input.Key,
		Value:     input.Value,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func GetConfigByID(db *sql.DB, id int64) (*models.Config, error) {
	config := &models.Config{}

	err := db.QueryRow(
		"SELECT id, key, value, created_at, updated_at FROM configs WHERE id = ?",
		id,
	).Scan(&config.ID, &config.Key, &config.Value, &config.CreatedAt, &config.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return config, nil
}

func GetConfigByKey(db *sql.DB, key string) (*models.Config, error) {
	config := &models.Config{}

	err := db.QueryRow(
		"SELECT id, key, value, created_at, updated_at FROM configs WHERE key = ?",
		key,
	).Scan(&config.ID, &config.Key, &config.Value, &config.CreatedAt, &config.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return config, nil
}

func GetAllConfigs(db *sql.DB) ([]*models.Config, error) {
	rows, err := db.Query(
		"SELECT id, key, value, created_at, updated_at FROM configs ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query configs: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var configs []*models.Config
	for rows.Next() {
		config := &models.Config{}
		if err := rows.Scan(&config.ID, &config.Key, &config.Value, &config.CreatedAt, &config.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func UpdateConfig(db *sql.DB, id int64, input models.UpdateConfigInput) (*models.Config, error) {
	config, err := GetConfigByID(db, id)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	config.Key = input.Key
	config.Value = input.Value
	config.UpdatedAt = time.Now()

	_, err = db.Exec(
		"UPDATE configs SET key = ?, value = ?, updated_at = ? WHERE id = ?",
		config.Key, config.Value, config.UpdatedAt, config.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	return config, nil
}

func DeleteConfig(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM configs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
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
