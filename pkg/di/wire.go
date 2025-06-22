//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	_ "github.com/lib/pq"
	"simple-product-api/pkg/logger"
	"simple-product-api/pkg/redis"

	httpHandler "simple-product-api/internal/product/delivery/http"
	"simple-product-api/internal/product/repository"
	"simple-product-api/internal/product/usecase"
	"simple-product-api/pkg/config"
)

func InitializeHandler(cfg *config.Config) (*httpHandler.Handler, error) {
	wire.Build(
		ProvidePostgres,

		repository.NewPostgresRepo,
		wire.Bind(new(repository.ProductRepository), new(*repository.RepositoryPostgre)),

		usecase.NewUsecase,
		wire.Bind(new(usecase.ProductUsecase), new(*usecase.Usecase)),

		httpHandler.NewHandler,
		redis.NewRedis,

		logger.NewLogger,
	)
	return &httpHandler.Handler{}, nil
}
