package order

import "time"

const StatusCreated = "CREATED"
const StatusShipped = "SHIPPED"

type ItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
type CreateOrderRequest struct {
	Items []ItemRequest `json:"items"`
}
type OrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
type Order struct {
	ID        int64       `json:"id"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	ShippedAt *time.Time  `json:"shipped_at,omitempty"`
	Items     []OrderItem `json:"items"`
}
type ShipResponse struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}
