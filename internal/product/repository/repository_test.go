package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"simple-product-api/internal/product"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"simple-product-api/internal/product/repository"
)

func TestRepo_FindByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	expected := &product.Product{
		ID: "84b6f675-1e28-4ef4-b987-2e7422b4f5a0", Name: "Banana", Type: "Buah", Price: 10000, CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "type", "price", "created_at"}).
		AddRow(expected.ID, expected.Name, expected.Type, expected.Price, expected.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\$1").
		WithArgs(expected.ID).
		WillReturnRows(rows)

	result, err := repo.FindByID(context.Background(), expected.ID)

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestRepo_FindByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	id := "non-existent-id"

	mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\$1").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.FindByID(context.Background(), id)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestRepo_FindByID_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	id := "error-id"

	mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\$1").
		WithArgs(id).
		WillReturnError(errors.New("db connection lost"))

	_, err := repo.FindByID(context.Background(), id)
	assert.EqualError(t, err, "db connection lost")
}

func TestRepo_Find_NoFilter(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	filter := product.ListFilter{Page: 1, PageSize: 10}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "type", "price", "created_at", "total_count"}).
		AddRow("1", "Apple", "Buah", 15000, now, 1)

	mock.ExpectQuery(`SELECT id, name, type, price, created_at, COUNT\(\*\) OVER\(\) as total_count FROM products ORDER BY created_at DESC LIMIT 10 OFFSET 0`).
		WillReturnRows(rows)

	products, total, err := repo.Find(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, products, 1)
	assert.Equal(t, "Apple", products[0].Name)
}

func TestRepo_Find_WithQueryAndType(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	filter := product.ListFilter{Page: 1, PageSize: 5, Query: "banana", Type: "Buah", SortBy: "name", Order: "asc"}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "type", "price", "created_at", "total_count"}).
		AddRow("2", "Banana", "Buah", 12000, now, 1)

	mock.ExpectQuery(`SELECT id, name, type, price, created_at, COUNT\(\*\) OVER\(\) as total_count FROM products WHERE LOWER\(name\) LIKE LOWER\(\$1\) AND type = \$2 ORDER BY name ASC LIMIT 5 OFFSET 0`).
		WithArgs("%banana%", "Buah").
		WillReturnRows(rows)

	products, total, err := repo.Find(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, "Banana", products[0].Name)
}

func TestRepo_Find_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	filter := product.ListFilter{Page: 1, PageSize: 10}

	mock.ExpectQuery("SELECT id, name, type, price").
		WillReturnError(errors.New("db query error"))

	_, _, err := repo.Find(context.Background(), filter)

	assert.Error(t, err)
	assert.EqualError(t, err, "db query error")
}

func TestRepo_Find_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	filter := product.ListFilter{Page: 1, PageSize: 10}

	// simulate broken row (wrong column count)
	rows := sqlmock.NewRows([]string{"id", "name", "type", "price", "created_at"}). // missing total_count
											AddRow("3", "Carrot", "Sayuran", 8000, time.Now())

	mock.ExpectQuery("SELECT id, name, type, price").
		WillReturnRows(rows)

	_, _, err := repo.Find(context.Background(), filter)
	assert.Error(t, err)
}

func TestRepo_Save_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	p := &product.Product{
		ID: "123", Name: "Mango", Type: "Buah", Price: 13000, CreatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO products").
		WithArgs(p.ID, p.Name, p.Type, p.Price, p.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Save(context.Background(), p)
	assert.NoError(t, err)
}

func TestRepo_Save_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := logrus.New()
	repo := repository.NewPostgresRepo(db, log)

	p := &product.Product{
		ID: "123", Name: "Papaya", Type: "Buah", Price: 11000, CreatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO products").
		WithArgs(p.ID, p.Name, p.Type, p.Price, p.CreatedAt).
		WillReturnError(errors.New("insert error"))

	err := repo.Save(context.Background(), p)
	assert.Error(t, err)
	assert.EqualError(t, err, "insert error")
}
