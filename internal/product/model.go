package product

import "time"

type Product struct {
	ID          int64     `json:"id"`
	ArticleCode string    `json:"articleCode"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateProductRequest struct {
	ArticleCode string `json:"articleCode"`
	Name        string `json:"name"`
}
