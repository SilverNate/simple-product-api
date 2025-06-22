package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func SeedProducts(db *sql.DB) error {
	products := []struct {
		Name  string
		Type  string
		Price float64
	}{
		{"Tomato", "Sayuran", 5000},
		{"Chicken Breast", "Protein", 25000},
		{"Apple", "Buah", 8000},
		{"Chips", "Snack", 6000},
	}

	for _, p := range products {
		_, err := db.ExecContext(context.Background(), `
			INSERT INTO products (id, name, type, price, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`, uuid.New().String(), p.Name, p.Type, p.Price, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
