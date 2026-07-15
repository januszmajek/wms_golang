package product

import "time"

type Product struct {
	ID        int64     `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	SKU  string `json:"sku"`
	Name string `json:"name"`
}
