// Package db provides database operations for the application.
package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/koban/ci/models"
)

// CreateProduct creates a new product in the database.
func CreateProduct(db *sql.DB, input models.CreateProductInput) (*models.Product, error) {
	now := time.Now()

	result, err := db.Exec(
		"INSERT INTO products (name, price, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		input.Name, input.Price, input.Description, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &models.Product{
		ID:          id,
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// GetProductByID retrieves a product by ID.
func GetProductByID(db *sql.DB, id int64) (*models.Product, error) {
	product := &models.Product{}

	err := db.QueryRow(
		"SELECT id, name, price, description, created_at, updated_at FROM products WHERE id = ?",
		id,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.CreatedAt, &product.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetAllProducts retrieves all products from the database.
func GetAllProducts(db *sql.DB) ([]*models.Product, error) {
	rows, err := db.Query(
		"SELECT id, name, price, description, created_at, updated_at FROM products ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct updates an existing product.
func UpdateProduct(db *sql.DB, id int64, input models.UpdateProductInput) (*models.Product, error) {
	product, err := GetProductByID(db, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	product.Name = input.Name
	product.Price = input.Price
	product.Description = input.Description
	product.UpdatedAt = time.Now()

	_, err = db.Exec(
		"UPDATE products SET name = ?, price = ?, description = ?, updated_at = ? WHERE id = ?",
		product.Name, product.Price, product.Description, product.UpdatedAt, product.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// DeleteProduct deletes a product by ID.
func DeleteProduct(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
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
