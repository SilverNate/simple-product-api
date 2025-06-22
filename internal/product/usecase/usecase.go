package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"simple-product-api/internal/product"
	"time"

	"github.com/google/uuid"
	"simple-product-api/internal/product/repository"
)

type Usecase struct {
	Repo  repository.ProductRepository
	Redis *redis.Client
	Log   *logrus.Logger
}

func NewUsecase(repo repository.ProductRepository, redis *redis.Client, log *logrus.Logger) *Usecase {
	return &Usecase{Repo: repo, Redis: redis, Log: log}
}

func (uc *Usecase) CreateProduct(ctx context.Context, product *product.Product) error {
	uc.Log.WithFields(logrus.Fields{
		"name":  product.Name,
		"type":  product.Type,
		"price": product.Price,
	}).Info("creating new product")

	existing, err := uc.Repo.FindProductByNameAndType(ctx, product.Name, product.Type)
	if err != nil {
		uc.Log.Error("failed to check existing product: ", err)
		return err
	}
	if existing != nil {
		return fmt.Errorf("product with name '%s' and type '%s' already exists", product.Name, product.Type)
	}

	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()

	err = uc.Repo.SaveProduct(ctx, product)
	if err != nil {
		uc.Log.Error("error save product: ", err)
		return err
	}

	return nil
}

func (uc *Usecase) ListProduct(ctx context.Context, filter product.ListFilter) (products []product.Product, total int, err error) {
	uc.Log.WithFields(logrus.Fields{
		"query": filter.Query,
		"type":  filter.Type,
		"page":  filter.Page,
		"size":  filter.PageSize,
	}).Info("listing products")

	cacheKey := fmt.Sprintf("products:all:name=%s:type=%s:sort=%s:order=%s:page=%d:size=%d",
		filter.Query, filter.Type, filter.SortBy, filter.Order, filter.Page, filter.PageSize,
	)

	cached, err := uc.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		uc.Log.WithField("cache_key", cacheKey).Info("redis cache present")
		var cachedProducts []product.Product
		if err := json.Unmarshal([]byte(cached), &cachedProducts); err == nil {
			return cachedProducts, len(cachedProducts), nil
		}
	} else if err == redis.Nil {
		uc.Log.WithField("cache_key", cacheKey).Warn("redis cache missing")
	} else {
		uc.Log.WithError(err).Error("redis error")
	}

	products, total, err = uc.Repo.FindProduct(ctx, filter)
	if err != nil {
		return products, total, err
	}

	data, err := json.Marshal(products)
	if err != nil {
		uc.Log.Errorf("failed to marshal products for caching: %v", err)
		return
	}

	go uc.Redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return products, total, nil
}

func (uc *Usecase) GetProductByID(ctx context.Context, id string) (products *product.Product, err error) {
	uc.Log.WithField("id", id).Info("retrieving product by ID")

	cacheKey := "products:id:" + id

	var cached string
	err = retry.Do(func() error {
		var err error
		cached, err = uc.Redis.Get(ctx, cacheKey).Result()
		return err
	}, retry.Attempts(3), retry.DelayType(retry.BackOffDelay))

	if cached != "" {
		if err := json.Unmarshal([]byte(cached), &products); err != nil {
			uc.Log.Errorf("error unmarshall from redis: %v", err.Error())
		} else {
			return products, nil
		}
	}

	products, err = uc.Repo.FindProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(products)
	if err != nil {
		uc.Log.Errorf("failed to marshal products for caching: %v", err)
		return
	}

	go uc.Redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return products, nil
}
