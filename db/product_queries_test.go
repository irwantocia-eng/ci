package db

import (
	"testing"

	"github.com/koban/ci/models"
)

func TestCreateProduct(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	product, err := CreateProduct(database.DB, models.CreateProductInput{
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test description",
	})
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if product.ID == 0 {
		t.Error("Expected product ID to be set")
	}

	if product.Name != "Test Product" {
		t.Errorf("Expected name 'Test Product', got '%s'", product.Name)
	}

	if product.Price != 99.99 {
		t.Errorf("Expected price 99.99, got %f", product.Price)
	}
}

func TestGetProductByID(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdProduct, err := CreateProduct(database.DB, models.CreateProductInput{
		Name:        "Get Test Product",
		Price:       49.99,
		Description: "Get test description",
	})
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	product, err := GetProductByID(database.DB, createdProduct.ID)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if product == nil {
		t.Fatal("Expected product to be found")
	}

	if product.ID != createdProduct.ID {
		t.Errorf("Expected ID %d, got %d", createdProduct.ID, product.ID)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	product, err := GetProductByID(database.DB, 999)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if product != nil {
		t.Error("Expected nil product for non-existent ID")
	}
}

func TestGetAllProducts(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	_, err := CreateProduct(database.DB, models.CreateProductInput{
		Name:  "Product 1",
		Price: 10.00,
	})
	if err != nil {
		t.Fatalf("Failed to create product 1: %v", err)
	}

	_, err = CreateProduct(database.DB, models.CreateProductInput{
		Name:  "Product 2",
		Price: 20.00,
	})
	if err != nil {
		t.Fatalf("Failed to create product 2: %v", err)
	}

	products, err := GetAllProducts(database.DB)
	if err != nil {
		t.Fatalf("Failed to get all products: %v", err)
	}

	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}
}

func TestUpdateProduct(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdProduct, err := CreateProduct(database.DB, models.CreateProductInput{
		Name:  "Original Product",
		Price: 50.00,
	})
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	updatedProduct, err := UpdateProduct(database.DB, createdProduct.ID, models.UpdateProductInput{
		Name:        "Updated Product",
		Price:       75.50,
		Description: "Updated description",
	})
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	if updatedProduct.Name != "Updated Product" {
		t.Errorf("Expected name 'Updated Product', got '%s'", updatedProduct.Name)
	}

	if updatedProduct.Price != 75.50 {
		t.Errorf("Expected price 75.50, got %f", updatedProduct.Price)
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	product, err := UpdateProduct(database.DB, 999, models.UpdateProductInput{
		Name:  "Test",
		Price: 10.00,
	})
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	if product != nil {
		t.Error("Expected nil product for non-existent ID")
	}
}

func TestDeleteProduct(t *testing.T) {
	database := setupTestDB(t)
	defer func() { _ = database.Close() }()

	createdProduct, err := CreateProduct(database.DB, models.CreateProductInput{
		Name:  "Delete Me",
		Price: 5.00,
	})
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	err = DeleteProduct(database.DB, createdProduct.ID)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}

	product, _ := GetProductByID(database.DB, createdProduct.ID)
	if product != nil {
		t.Error("Expected product to be deleted")
	}
}
