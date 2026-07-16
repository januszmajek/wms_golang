package order

import "errors"

var (
	ErrEmptyOrder        = errors.New("order needs at least one item")
	ErrBadItemQuantity   = errors.New("item quantity must be bigger than 0")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrAlreadyShipped    = errors.New("order already shipped")
)

// Store is the only production interface in the application. The order service
// uses several storage operations, so this boundary keeps its business rules
// testable without a database.
type Store interface {
	GetStock(productID int64) (int, error)
	Create(items []OrderItem) (Order, error)
	Get(id int64) (Order, error)
	Ship(id int64) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) Create(req CreateOrderRequest) (Order, error) {
	if len(req.Items) == 0 {
		return Order{}, ErrEmptyOrder
	}
	merged := map[int64]int{}
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return Order{}, ErrBadItemQuantity
		}
		merged[item.ProductID] += item.Quantity
	}
	var items []OrderItem
	for productID, qty := range merged {
		stock, err := s.store.GetStock(productID)
		if err != nil {
			return Order{}, err
		}
		if stock < qty {
			return Order{}, ErrInsufficientStock
		}
		items = append(items, OrderItem{ProductID: productID, Quantity: qty})
	}
	return s.store.Create(items)
}
func (s *Service) Get(id int64) (Order, error) { return s.store.Get(id) }
func (s *Service) Ship(id int64) (ShipResponse, error) {
	if err := s.store.Ship(id); err != nil {
		return ShipResponse{}, err
	}
	return ShipResponse{OrderID: id, Status: StatusShipped}, nil
}
