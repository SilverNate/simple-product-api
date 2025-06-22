package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"simple-product-api/internal/product"
	"strings"
)

type RepositoryPostgre struct {
	db  *sql.DB
	Log *logrus.Logger
}

func NewPostgresRepo(db *sql.DB, log *logrus.Logger) *RepositoryPostgre {
	return &RepositoryPostgre{db: db, Log: log}
}

func (r *RepositoryPostgre) SaveProduct(ctx context.Context, p *product.Product) error {
	query := `INSERT INTO products (id, name, type, price, created_at)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Type, p.Price, p.CreatedAt)
	if err != nil {
		r.Log.WithError(err).Error("error inserting product")
		return err
	}
	return err
}

func (r *RepositoryPostgre) FindProductByID(ctx context.Context, id string) (*product.Product, error) {
	query := `SELECT id, name, type, price, created_at FROM products WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var p product.Product
	err := row.Scan(&p.ID, &p.Name, &p.Type, &p.Price, &p.CreatedAt)
	if err != nil {
		r.Log.WithError(err).Errorf("error find product by id: %v", id)
		return nil, err
	}
	return &p, nil
}

func (r *RepositoryPostgre) FindProduct(ctx context.Context, f product.ListFilter) (products []product.Product, total int, err error) {
	baseQuery := `SELECT id, name, type, price, created_at, COUNT(*) OVER() as total_count FROM products`
	var clauses []string
	var args []interface{}

	argIndex := 1

	if f.Query != "" {
		clauses = append(clauses, fmt.Sprintf("LOWER(name) LIKE LOWER($%d)", argIndex))
		args = append(args, "%"+f.Query+"%")
		argIndex++
	}

	if f.Type != "" {
		clauses = append(clauses, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, f.Type)
		argIndex++
	}

	if len(clauses) > 0 {
		baseQuery += " WHERE " + strings.Join(clauses, " AND ")
	}

	orderBy := "created_at"
	if f.SortBy == "name" || f.SortBy == "price" || f.SortBy == "created_at" {
		orderBy = f.SortBy
	}
	order := strings.ToUpper(f.Order)
	if order != "ASC" {
		order = "DESC"
	}
	baseQuery += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)

	limit := f.PageSize
	offset := (f.Page - 1) * f.PageSize
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		r.Log.WithError(err).Error(fmt.Sprintf("error find product using filter: %+v", f))
		return nil, total, err
	}
	defer rows.Close()

	for rows.Next() {
		var p product.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Price, &p.CreatedAt, &total); err != nil {
			r.Log.WithError(err).Error("error row scan in find product using filter")
			return nil, total, err
		}
		products = append(products, p)
	}
	return products, total, nil
}

func (r *RepositoryPostgre) FindProductByNameAndType(ctx context.Context, name, ptype string) (*product.Product, error) {
	query := `SELECT id, name, type, price, created_at FROM products WHERE LOWER(name) = LOWER($1) AND LOWER(type) = LOWER($2) LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, name, ptype)

	var product product.Product
	err := row.Scan(&product.ID, &product.Name, &product.Type, &product.Price, &product.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.Log.WithError(err).Error("error find product same name and type")
		return nil, err
	}
	return &product, nil
}
