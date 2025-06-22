package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"simple-product-api/pkg/config"
)

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		logrus.Fatalf("migration failed: %v", err)
	}

	if err := SeedProducts(db); err != nil {
		logrus.Fatalf("seeding failed: %v", err)
	}

	return db, nil
}
