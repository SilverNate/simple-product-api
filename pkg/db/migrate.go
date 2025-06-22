package db

import (
	"database/sql"
	"os"
)

func Migrate(db *sql.DB) error {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM information_schema.tables 
		WHERE table_name = 'products'
	)`).Scan(&exists)

	if err != nil || exists {
		return nil
	}

	script, err := os.ReadFile("migrations/001_create_products.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(script))
	return err
}
