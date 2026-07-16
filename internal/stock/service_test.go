package stock

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestService(t *testing.T) (*Service, sqlmock.Sqlmock) {
	t.Helper()

	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet SQL expectation: %v", err)
		}
		database.Close()
	})

	return NewService(NewRepository(database)), mock
}

func TestReceiveAddsStock(t *testing.T) {
	service, mock := newTestService(t)
	mock.ExpectQuery("SELECT EXISTS").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO stock").WithArgs(int64(1), 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO inbound_operations").WithArgs(int64(1), 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	response, err := service.Receive(InboundRequest{ProductID: 1, Quantity: 10})
	if err != nil {
		t.Fatalf("receive stock: %v", err)
	}
	if response != (InboundResponse{ProductID: 1, QuantityAdded: 10}) {
		t.Fatalf("unexpected response: %#v", response)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestReceiveRejectsBadQuantity(t *testing.T) {
	service, _ := newTestService(t)

	_, err := service.Receive(InboundRequest{ProductID: 1, Quantity: 0})
	if !errors.Is(err, ErrBadQuantity) {
		t.Fatalf("wanted bad quantity, got %v", err)
	}
}

func TestReceiveRejectsMissingProduct(t *testing.T) {
	service, mock := newTestService(t)
	mock.ExpectQuery("SELECT EXISTS").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	_, err := service.Receive(InboundRequest{ProductID: 1, Quantity: 1})
	if !errors.Is(err, ErrProductMissing) {
		t.Fatalf("wanted missing product, got %v", err)
	}
}

func TestReceiveReturnsLookupError(t *testing.T) {
	service, mock := newTestService(t)
	wantErr := errors.New("database unavailable")
	mock.ExpectQuery("SELECT EXISTS").WithArgs(int64(1)).WillReturnError(wantErr)

	_, err := service.Receive(InboundRequest{ProductID: 1, Quantity: 1})
	if !errors.Is(err, wantErr) {
		t.Fatalf("wanted lookup error, got %v", err)
	}
}

func TestReceiveReturnsInboundError(t *testing.T) {
	service, mock := newTestService(t)
	wantErr := errors.New("cannot start transaction")
	mock.ExpectQuery("SELECT EXISTS").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectBegin().WillReturnError(wantErr)

	_, err := service.Receive(InboundRequest{ProductID: 1, Quantity: 1})
	if !errors.Is(err, wantErr) {
		t.Fatalf("wanted inbound error, got %v", err)
	}
}

func TestReportReturnsRepositoryItems(t *testing.T) {
	service, mock := newTestService(t)
	mock.ExpectQuery("SELECT p.id").WillReturnRows(
		sqlmock.NewRows([]string{"id", "article_code", "name", "quantity"}).
			AddRow(1, "CHAIR-1", "Chair", 8),
	)

	items, err := service.Report()
	if err != nil {
		t.Fatalf("get report: %v", err)
	}
	if len(items) != 1 || items[0].Quantity != 8 {
		t.Fatalf("unexpected report: %#v", items)
	}
}
