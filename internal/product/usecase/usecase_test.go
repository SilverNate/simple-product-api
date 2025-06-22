package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	redismock "github.com/go-redis/redismock/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"simple-product-api/internal/product"
	//mockRepo "simple-product-api/internal/product/mocks"
	mockRepo "simple-product-api/internal/product/mocks"
	"simple-product-api/internal/product/usecase"
	"testing"
	"time"
)

type UsecaseProductTestSuite struct {
	suite.Suite
	usecase   *usecase.Usecase
	mockRepo  *mockRepo.ProductRepository
	redisMock redismock.ClientMock
}

func (s *UsecaseProductTestSuite) SetupTest() {
	rdb, mock := redismock.NewClientMock()
	s.redisMock = mock
	s.mockRepo = mockRepo.NewProductRepository(s.T())
	logger := logrus.New()
	s.usecase = usecase.NewUsecase(s.mockRepo, rdb, logger)
}

func (s *UsecaseProductTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
}

func TestUsecaseProductTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseProductTestSuite))
}

func (s *UsecaseProductTestSuite) TestGetByIDWithRedisPresentSuccess() {
	expectedProduct := &product.Product{
		ID:    "123",
		Name:  "Sawi",
		Type:  "Sayuran",
		Price: float64(5000),
	}

	cacheKey := "products:id:123"
	data, _ := json.Marshal(expectedProduct)

	s.redisMock.ExpectGet(cacheKey).SetVal(string(data))

	// just in case
	s.mockRepo.AssertNotCalled(s.T(), "FindProductByID", mock.Anything, "123")

	res, err := s.usecase.GetProductByID(context.Background(), "123")

	s.NoError(err)
	s.Equal(expectedProduct.Name, res.Name)
	s.Equal(expectedProduct.Type, res.Type)
	s.Equal(expectedProduct.Price, res.Price)

	err = s.redisMock.ExpectationsWereMet()
	s.NoError(err)
}

func (s *UsecaseProductTestSuite) TestCreateSuccess() {
	s.mockRepo.On("FindProductByNameAndType", mock.Anything, "Banana", "Buah").Return(nil, nil)

	s.mockRepo.On("SaveProduct", mock.Anything, mock.Anything).Return(nil)

	newProduct := &product.Product{
		Name:  "Banana",
		Type:  "Buah",
		Price: 10000,
	}

	err := s.usecase.CreateProduct(context.Background(), newProduct)
	s.NoError(err)
}

func (s *UsecaseProductTestSuite) TestListProductWithRedisEmptySuccess() {
	products := []product.Product{
		{ID: "1", Name: "A", Type: "Buah", Price: 10000, CreatedAt: time.Now()},
	}
	filter := product.ListFilter{Page: 1, PageSize: 10}
	cacheKey := "products:all:name=:type=:sort=:order=:page=1:size=10"

	s.redisMock.ExpectGet(cacheKey).RedisNil()
	s.redisMock.ExpectSet(cacheKey, mock.Anything, 5*time.Minute).SetVal("OK")

	s.mockRepo.On("FindProduct", mock.Anything, mock.Anything).Return(products, 1, nil)

	res, total, err := s.usecase.ListProduct(context.Background(), filter)

	s.NoError(err)
	s.Equal(1, total)
	s.Equal("A", res[0].Name)
}

func (s *UsecaseProductTestSuite) TestListProductWithRedisPresentSuccess() {

	products := []product.Product{
		{ID: "1", Name: "A", Type: "Buah", Price: 10000, CreatedAt: time.Now()},
	}
	filter := product.ListFilter{Page: 1, PageSize: 10}
	cacheKey := "products:all:name=:type=:sort=:order=:page=1:size=10"

	jsonData, _ := json.Marshal(products)
	s.redisMock.ExpectGet(cacheKey).SetVal(string(jsonData))

	res, total, err := s.usecase.ListProduct(context.Background(), filter)

	s.NoError(err)
	s.Equal(1, total)
	s.Len(res, 1)
	s.Equal("A", res[0].Name)
}

func (s *UsecaseProductTestSuite) TestListProductWithRedisError() {

	products := []product.Product{
		{ID: "1", Name: "A", Type: "Buah", Price: 10000, CreatedAt: time.Now()},
	}
	filter := product.ListFilter{Page: 1, PageSize: 10}
	cacheKey := fmt.Sprintf("products:all:name=%s:type=%s:sort=%s:order=%s:page=%d:size=%d",
		filter.Query, filter.Type, filter.SortBy, filter.Order, filter.Page, filter.PageSize,
	)

	s.redisMock.ExpectGet(cacheKey).SetErr(errors.New("simulated redis connection error"))
	s.redisMock.ExpectSet(cacheKey, mock.Anything, 5*time.Minute).SetVal("OK")

	s.mockRepo.On("FindProduct", mock.Anything, mock.Anything).Return(products, 1, nil).Once()

	res, total, err := s.usecase.ListProduct(context.Background(), filter)

	s.NoError(err)
	s.Equal(1, total)
	s.Len(res, 1)
	s.Equal("A", res[0].Name)

}
