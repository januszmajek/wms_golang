package product

import "time"

type Product struct {
	ID          int64     `json:"id"`
	ArticleCode string    `json:"article_code"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	ArticleCode string `json:"article_code"`
	Name        string `json:"name"`
}
