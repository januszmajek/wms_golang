package order

import "errors"

var ErrEmptyOrder = errors.New("order needs at least one item")
var ErrBadItemQuantity = errors.New("item quantity must be bigger than 0")
var ErrInsufficientStock = errors.New("insufficient stock")
var ErrAlreadyShipped = errors.New("order already shipped")

type Service struct{ Repo RepositoryInterface }

func NewService(repo RepositoryInterface) *Service { return &Service{Repo: repo} }

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
		stock, err := s.Repo.GetStock(productID)
		if err != nil {
			return Order{}, err
		}
		if stock < qty {
			return Order{}, ErrInsufficientStock
		}
		items = append(items, OrderItem{ProductID: productID, Quantity: qty})
	}
	return s.Repo.Create(items)
}
func (s *Service) Get(id int64) (Order, error) { return s.Repo.Get(id) }
func (s *Service) Ship(id int64) (ShipResponse, error) {
	if err := s.Repo.Ship(id); err != nil {
		return ShipResponse{}, err
	}
	return ShipResponse{OrderID: id, Status: StatusShipped}, nil
}
