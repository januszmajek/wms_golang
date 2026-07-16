package order

import "time"

const StatusCreated = "CREATED"
const StatusShipped = "SHIPPED"

type ItemRequest struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}
type CreateOrderRequest struct {
	Items []ItemRequest `json:"items"`
}
type OrderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}
type Order struct {
	ID        int64       `json:"id"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	ShippedAt *time.Time  `json:"shippedAt,omitempty"`
	Items     []OrderItem `json:"items"`
}
type ShipResponse struct {
	OrderID int64  `json:"orderId"`
	Status  string `json:"status"`
}
