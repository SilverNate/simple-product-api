package product

import "time"

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=3"`
	Type      string    `json:"type" validate:"required,oneof=Sayuran Protein Buah Snack"`
	Price     float64   `json:"price" validate:"required,gt=0"`
	CreatedAt time.Time `json:"created_at"`
}

type ListFilter struct {
	Query    string
	Type     string
	SortBy   string
	Order    string
	Page     int
	PageSize int
}

type PaginatedResult struct {
	Items interface{} `json:"items"`
	Meta  MetaPage    `json:"meta"`
}

type MetaPage struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
