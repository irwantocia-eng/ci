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

func TestConfigListHandler_GET(t *testing.T) {
	database := setupTestDB(t)
	defer func() {
		_ = database.Close()
	}()

	_, err := db.CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "app_name",
		Value: "Test App",
	})
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/configs", nil)
	w := httptest.NewRecorder()

	ConfigListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var configs []models.Config
	if err := json.NewDecoder(res.Body).Decode(&configs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(configs) != 1 {
		t.Errorf("Expected 1 config, got %d", len(configs))
	}

	if configs[0].Key != "app_name" {
		t.Errorf("Expected key 'app_name', got '%s'", configs[0].Key)
	}
}

func TestConfigListHandler_POST(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	configInput := models.CreateConfigInput{
		Key:   "new_key",
		Value: "new_value",
	}
	body, _ := json.Marshal(configInput)

	req := httptest.NewRequest(http.MethodPost, "/api/configs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ConfigListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", res.StatusCode)
	}

	var config models.Config
	if err := json.NewDecoder(res.Body).Decode(&config); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if config.Key != configInput.Key {
		t.Errorf("Expected key '%s', got '%s'", configInput.Key, config.Key)
	}

	if config.Value != configInput.Value {
		t.Errorf("Expected value '%s', got '%s'", configInput.Value, config.Value)
	}
}

func TestConfigHandler_GET(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdConfig, err := db.CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "test_key",
		Value: "test_value",
	})
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/configs/1", nil)
	w := httptest.NewRecorder()

	ConfigHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var config models.Config
	if err := json.NewDecoder(res.Body).Decode(&config); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if config.ID != createdConfig.ID {
		t.Errorf("Expected config ID %d, got %d", createdConfig.ID, config.ID)
	}
}

func TestConfigHandler_GET_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodGet, "/api/configs/999", nil)
	w := httptest.NewRecorder()

	ConfigHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", res.StatusCode)
	}
}

func TestConfigHandler_PUT(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := db.CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "original_key",
		Value: "original_value",
	})
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	updateInput := models.UpdateConfigInput{
		Key:   "updated_key",
		Value: "updated_value",
	}
	body, _ := json.Marshal(updateInput)

	req := httptest.NewRequest(http.MethodPut, "/api/configs/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ConfigHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var config models.Config
	if err := json.NewDecoder(res.Body).Decode(&config); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if config.Key != updateInput.Key {
		t.Errorf("Expected key '%s', got '%s'", updateInput.Key, config.Key)
	}

	if config.Value != updateInput.Value {
		t.Errorf("Expected value '%s', got '%s'", updateInput.Value, config.Value)
	}
}

func TestConfigHandler_DELETE(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := db.CreateConfig(database.DB, models.CreateConfigInput{
		Key:   "delete_key",
		Value: "delete_value",
	})
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/configs/1", nil)
	w := httptest.NewRecorder()

	ConfigHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", res.StatusCode)
	}

	config, err := db.GetConfigByID(database.DB, 1)
	if err != nil {
		t.Fatalf("Failed to check config: %v", err)
	}

	if config != nil {
		t.Error("Expected config to be deleted, but config still exists")
	}
}

func TestConfigListHandler_MethodNotAllowed(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodPatch, "/api/configs", nil)
	w := httptest.NewRecorder()

	ConfigListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", res.StatusCode)
	}
}
