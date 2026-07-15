package stock

import "testing"

type fakeStockRepo struct {
	exists    bool
	addCalled bool
}

func (f *fakeStockRepo) ProductExists(productID int64) (bool, error) { return f.exists, nil }
func (f *fakeStockRepo) AddInbound(productID int64, quantity int) error {
	f.addCalled = true
	return nil
}
func (f *fakeStockRepo) Report() ([]ReportItem, error) { return nil, nil }

func TestReceiveAddsStock(t *testing.T) {
	repo := &fakeStockRepo{exists: true}
	service := NewService(repo)
	resp, err := service.Receive(InboundRequest{ProductID: 1, Quantity: 10})
	if err != nil {
		t.Fatal(err)
	}
	if resp.QuantityAdded != 10 || !repo.addCalled {
		t.Fatalf("stock was not added")
	}
}
func TestReceiveRejectsBadQuantity(t *testing.T) {
	service := NewService(&fakeStockRepo{exists: true})
	_, err := service.Receive(InboundRequest{ProductID: 1, Quantity: 0})
	if err != ErrBadQuantity {
		t.Fatalf("wanted bad quantity got %v", err)
	}
}
func TestReceiveRejectsMissingProduct(t *testing.T) {
	service := NewService(&fakeStockRepo{exists: false})
	_, err := service.Receive(InboundRequest{ProductID: 1, Quantity: 1})
	if err != ErrProductMissing {
		t.Fatalf("wanted missing product got %v", err)
	}
}
