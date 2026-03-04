package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/koban/ci/db"
	"github.com/koban/ci/models"
)

func ConfigListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllConfigs(w)
	case http.MethodPost:
		createConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/api/configs/")
	if err != nil {
		http.Error(w, "Invalid config ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getConfigByID(w, id)
	case http.MethodPut:
		updateConfig(w, r, id)
	case http.MethodDelete:
		deleteConfig(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllConfigs(w http.ResponseWriter) {
	configs, err := db.GetAllConfigs(dbInstance.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(configs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getConfigByID(w http.ResponseWriter, id int64) {
	config, err := db.GetConfigByID(dbInstance.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if config == nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createConfig(w http.ResponseWriter, r *http.Request) {
	var input models.CreateConfigInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existing, err := db.GetConfigByKey(dbInstance.DB, input.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, "Config key already exists", http.StatusConflict)
		return
	}

	config, err := db.CreateConfig(dbInstance.DB, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateConfig(w http.ResponseWriter, r *http.Request, id int64) {
	var input models.UpdateConfigInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existing, err := db.GetConfigByKey(dbInstance.DB, input.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing != nil && existing.ID != id {
		http.Error(w, "Config key already exists", http.StatusConflict)
		return
	}

	config, err := db.UpdateConfig(dbInstance.DB, id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if config == nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteConfig(w http.ResponseWriter, id int64) {
	if err := db.DeleteConfig(dbInstance.DB, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
