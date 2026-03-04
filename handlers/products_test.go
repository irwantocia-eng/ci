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

func TestProductListHandler_GET(t *testing.T) {
	database := setupTestDB(t)
	defer func() {
		_ = database.Close()
	}()

	_, err := db.CreateProduct(database.DB, models.CreateProductInput{
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test description",
	})
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	w := httptest.NewRecorder()

	ProductListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var products []models.Product
	if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(products))
	}

	if products[0].Name != "Test Product" {
		t.Errorf("Expected name 'Test Product', got '%s'", products[0].Name)
	}
}

func TestProductListHandler_POST(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	productInput := models.CreateProductInput{
		Name:        "New Product",
		Price:       149.99,
		Description: "New product description",
	}
	body, _ := json.Marshal(productInput)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ProductListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", res.StatusCode)
	}

	var product models.Product
	if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if product.Name != productInput.Name {
		t.Errorf("Expected name '%s', got '%s'", productInput.Name, product.Name)
	}

	if product.Price != productInput.Price {
		t.Errorf("Expected price %.2f, got %.2f", productInput.Price, product.Price)
	}
}

func TestProductHandler_GET(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdProduct, err := db.CreateProduct(database.DB, models.CreateProductInput{
		Name:        "Get Product",
		Price:       79.99,
		Description: "Product to get",
	})
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	w := httptest.NewRecorder()

	ProductHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var product models.Product
	if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if product.ID != createdProduct.ID {
		t.Errorf("Expected product ID %d, got %d", createdProduct.ID, product.ID)
	}
}

func TestProductHandler_GET_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodGet, "/api/products/999", nil)
	w := httptest.NewRecorder()

	ProductHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", res.StatusCode)
	}
}

func TestProductHandler_PUT(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := db.CreateProduct(database.DB, models.CreateProductInput{
		Name:        "Original Product",
		Price:       50.00,
		Description: "Original",
	})
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	updateInput := models.UpdateProductInput{
		Name:        "Updated Product",
		Price:       75.50,
		Description: "Updated description",
	}
	body, _ := json.Marshal(updateInput)

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ProductHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	var product models.Product
	if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if product.Name != updateInput.Name {
		t.Errorf("Expected name '%s', got '%s'", updateInput.Name, product.Name)
	}

	if product.Price != updateInput.Price {
		t.Errorf("Expected price %.2f, got %.2f", updateInput.Price, product.Price)
	}
}

func TestProductHandler_DELETE(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := db.CreateProduct(database.DB, models.CreateProductInput{
		Name:        "Delete Product",
		Price:       25.00,
		Description: "Product to delete",
	})
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)
	w := httptest.NewRecorder()

	ProductHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", res.StatusCode)
	}

	// Verify product was deleted
	product, err := db.GetProductByID(database.DB, 1)
	if err != nil {
		t.Fatalf("Failed to check product: %v", err)
	}

	if product != nil {
		t.Error("Expected product to be deleted, but product still exists")
	}
}

func TestProductListHandler_MethodNotAllowed(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	req := httptest.NewRequest(http.MethodPatch, "/api/products", nil)
	w := httptest.NewRecorder()

	ProductListHandler(w, req)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", res.StatusCode)
	}
}
