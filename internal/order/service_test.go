package order

import (
	"testing"
	"time"
)

type fakeOrderRepo struct {
	stock        map[int64]int
	orders       map[int64]Order
	shipped      map[int64]bool
	createdItems []OrderItem
}

func newFakeOrderRepo() *fakeOrderRepo {
	return &fakeOrderRepo{stock: map[int64]int{}, orders: map[int64]Order{}, shipped: map[int64]bool{}}
}
func (f *fakeOrderRepo) GetStock(productID int64) (int, error) { return f.stock[productID], nil }
func (f *fakeOrderRepo) Create(items []OrderItem) (Order, error) {
	f.createdItems = items
	o := Order{ID: 1, Status: StatusCreated, CreatedAt: time.Now(), Items: items}
	f.orders[1] = o
	return o, nil
}
func (f *fakeOrderRepo) Get(id int64) (Order, error) { return f.orders[id], nil }
func (f *fakeOrderRepo) Ship(id int64) error {
	if f.shipped[id] {
		return ErrAlreadyShipped
	}
	o := f.orders[id]
	for _, item := range o.Items {
		if f.stock[item.ProductID] < item.Quantity {
			return ErrInsufficientStock
		}
	}
	for _, item := range o.Items {
		f.stock[item.ProductID] -= item.Quantity
	}
	f.shipped[id] = true
	return nil
}

func TestCreateOrderSuccess(t *testing.T) {
	repo := newFakeOrderRepo()
	repo.stock[1] = 10
	service := NewService(repo)
	o, err := service.Create(CreateOrderRequest{Items: []ItemRequest{{ProductID: 1, Quantity: 4}}})
	if err != nil {
		t.Fatal(err)
	}
	if o.Status != StatusCreated {
		t.Fatalf("bad status")
	}
}
func TestCreateOrderRejectsEmpty(t *testing.T) {
	_, err := NewService(newFakeOrderRepo()).Create(CreateOrderRequest{})
	if err != ErrEmptyOrder {
		t.Fatalf("wanted empty err got %v", err)
	}
}
func TestCreateOrderRejectsBadQuantity(t *testing.T) {
	_, err := NewService(newFakeOrderRepo()).Create(CreateOrderRequest{Items: []ItemRequest{{ProductID: 1, Quantity: -1}}})
	if err != ErrBadItemQuantity {
		t.Fatalf("wanted bad quantity got %v", err)
	}
}
func TestCreateOrderRejectsInsufficientStock(t *testing.T) {
	repo := newFakeOrderRepo()
	repo.stock[1] = 2
	_, err := NewService(repo).Create(CreateOrderRequest{Items: []ItemRequest{{ProductID: 1, Quantity: 3}}})
	if err != ErrInsufficientStock {
		t.Fatalf("wanted insufficient stock got %v", err)
	}
}
func TestCreateOrderSumsDuplicateItems(t *testing.T) {
	repo := newFakeOrderRepo()
	repo.stock[1] = 12
	_, err := NewService(repo).Create(CreateOrderRequest{Items: []ItemRequest{{ProductID: 1, Quantity: 6}, {ProductID: 1, Quantity: 6}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.createdItems) != 1 || repo.createdItems[0].Quantity != 12 {
		t.Fatalf("duplicates not summed: %#v", repo.createdItems)
	}
}
func TestCreateOrderMultipleProducts(t *testing.T) {
	repo := newFakeOrderRepo()
	repo.stock[1] = 5
	repo.stock[2] = 7
	_, err := NewService(repo).Create(CreateOrderRequest{Items: []ItemRequest{{ProductID: 1, Quantity: 5}, {ProductID: 2, Quantity: 7}}})
	if err != nil {
		t.Fatal(err)
	}
}
func TestShipOrderDecreasesStock(t *testing.T) {
	repo := newFakeOrderRepo()
	repo.stock[1] = 10
	service := NewService(repo)
	_, _ = service.Create(CreateOrderRequest{Items: []ItemRequest{{ProductID: 1, Quantity: 4}}})
	resp, err := service.Ship(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != StatusShipped || repo.stock[1] != 6 {
		t.Fatalf("ship failed")
	}
}
func TestShipOrderRejectsDoubleShip(t *testing.T) {
	repo := newFakeOrderRepo()
	repo.stock[1] = 10
	service := NewService(repo)
	_, _ = service.Create(CreateOrderRequest{Items: []ItemRequest{{ProductID: 1, Quantity: 4}}})
	_, _ = service.Ship(1)
	_, err := service.Ship(1)
	if err != ErrAlreadyShipped {
		t.Fatalf("wanted already shipped got %v", err)
	}
}
func TestShipOrderFailsIfStockDropped(t *testing.T) {
	repo := newFakeOrderRepo()
	repo.stock[1] = 10
	service := NewService(repo)
	_, _ = service.Create(CreateOrderRequest{Items: []ItemRequest{{ProductID: 1, Quantity: 4}}})
	repo.stock[1] = 1
	_, err := service.Ship(1)
	if err != ErrInsufficientStock {
		t.Fatalf("wanted insufficient stock got %v", err)
	}
}
