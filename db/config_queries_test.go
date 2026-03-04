package db

import (
	"testing"

	"github.com/koban/ci/models"
)

func TestCreateConfig(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	config, err := CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "app_name",
		Value: "Test App",
	})
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	if config.ID == 0 {
		t.Error("Expected config ID to be set")
	}

	if config.Key != "app_name" {
		t.Errorf("Expected key 'app_name', got '%s'", config.Key)
	}

	if config.Value != "Test App" {
		t.Errorf("Expected value 'Test App', got '%s'", config.Value)
	}
}

func TestGetConfigByID(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdConfig, err := CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "test_key",
		Value: "test_value",
	})
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	config, err := GetConfigByID(database.DB, createdConfig.ID)
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config == nil {
		t.Fatal("Expected config to be found")
	}

	if config.ID != createdConfig.ID {
		t.Errorf("Expected ID %d, got %d", createdConfig.ID, config.ID)
	}
}

func TestGetConfigByID_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	config, err := GetConfigByID(database.DB, 999)
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config != nil {
		t.Error("Expected nil config for non-existent ID")
	}
}

func TestGetConfigByKey(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "unique_key",
		Value: "unique_value",
	})
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	config, err := GetConfigByKey(database.DB, "unique_key")
	if err != nil {
		t.Fatalf("Failed to get config by key: %v", err)
	}

	if config == nil {
		t.Fatal("Expected config to be found")
	}

	if config.Value != "unique_value" {
		t.Errorf("Expected value 'unique_value', got '%s'", config.Value)
	}
}

func TestGetConfigByKey_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	config, err := GetConfigByKey(database.DB, "nonexistent")
	if err != nil {
		t.Fatalf("Failed to get config by key: %v", err)
	}

	if config != nil {
		t.Error("Expected nil config for non-existent key")
	}
}

func TestGetAllConfigs(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "key1",
		Value: "value1",
	})
	if err != nil {
		t.Fatalf("Failed to create config 1: %v", err)
	}

	_, err = CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "key2",
		Value: "value2",
	})
	if err != nil {
		t.Fatalf("Failed to create config 2: %v", err)
	}

	configs, err := GetAllConfigs(database.DB)
	if err != nil {
		t.Fatalf("Failed to get all configs: %v", err)
	}

	if len(configs) != 2 {
		t.Errorf("Expected 2 configs, got %d", len(configs))
	}
}

func TestUpdateConfig(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdConfig, err := CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "original_key",
		Value: "original_value",
	})
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	updatedConfig, err := UpdateConfig(database.DB, createdConfig.ID, models.UpdateConfigInput{
		Key:   "updated_key",
		Value: "updated_value",
	})
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	if updatedConfig.Key != "updated_key" {
		t.Errorf("Expected key 'updated_key', got '%s'", updatedConfig.Key)
	}

	if updatedConfig.Value != "updated_value" {
		t.Errorf("Expected value 'updated_value', got '%s'", updatedConfig.Value)
	}
}

func TestUpdateConfig_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	config, err := UpdateConfig(database.DB, 999, models.UpdateConfigInput{
		Key:   "test",
		Value: "test",
	})
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	if config != nil {
		t.Error("Expected nil config for non-existent ID")
	}
}

func TestDeleteConfig(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdConfig, err := CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "delete_key",
		Value: "delete_value",
	})
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	err = DeleteConfig(database.DB, createdConfig.ID)
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}

	config, _ := GetConfigByID(database.DB, createdConfig.ID)
	if config != nil {
		t.Error("Expected config to be deleted")
	}
}
