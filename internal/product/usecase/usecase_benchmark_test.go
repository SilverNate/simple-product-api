package usecase_test

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"simple-product-api/internal/product"
	mockRepo "simple-product-api/internal/product/mocks"
	"simple-product-api/internal/product/usecase"
	"testing"
	"time"
)

type benchmarkEnv struct {
	usecase   *usecase.Usecase
	mockRepo  *mockRepo.ProductRepository
	redisMock redismock.ClientMock
}

func setupBenchmarkEnv(tb testing.TB) *benchmarkEnv {
	rdb, mock := redismock.NewClientMock()
	logger := logrus.New()
	repo := mockRepo.NewProductRepository(tb)

	uc := usecase.NewUsecase(repo, rdb, logger)

	return &benchmarkEnv{
		usecase:   uc,
		mockRepo:  repo,
		redisMock: mock,
	}
}

func BenchmarkProductUsecase_GetProductByID(b *testing.B) {
	env := setupBenchmarkEnv(b)

	productID := "test-123"
	p := &product.Product{
		ID:        productID,
		Name:      "Tomato",
		Type:      "Sayuran",
		Price:     8000,
		CreatedAt: time.Now(),
	}

	env.mockRepo.On("FindProductByID", mock.Anything, productID).Return(p, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = env.usecase.GetProductByID(context.Background(), productID)
	}
}

func BenchmarkProductUsecase_ListProduct(b *testing.B) {
	env := setupBenchmarkEnv(b)

	products := []product.Product{
		{ID: "1", Name: "A", Type: "Buah", Price: 10000, CreatedAt: time.Now()},
	}

	filter := product.ListFilter{Page: 1, PageSize: 10}

	env.mockRepo.On("FindProduct", mock.Anything, mock.AnythingOfType("product.ListFilter")).Return(products, 1, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = env.usecase.ListProduct(context.Background(), filter)
	}
}

func BenchmarkProductUsecase_CreateProduct(b *testing.B) {
	env := setupBenchmarkEnv(b)
	env.mockRepo.On("FindProductByNameAndType", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	env.mockRepo.On("SaveProduct", mock.Anything, mock.Anything).Return(nil)

	p := &product.Product{
		ID:        "abc123",
		Name:      "Sapi",
		Type:      "Protein",
		Price:     45000,
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = env.usecase.CreateProduct(context.Background(), p)
	}
}
