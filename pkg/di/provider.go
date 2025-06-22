package di

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"simple-product-api/pkg/config"
	"simple-product-api/pkg/db"
)

func ProvidePostgres(cfg *config.Config, log *logrus.Logger) (*sql.DB, error) {
	log.Infof("Connecting to PostgreSQL: %s", cfg.PostgresDSN)
	return db.NewPostgres(cfg)
}
