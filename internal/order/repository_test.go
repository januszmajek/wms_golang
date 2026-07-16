package order

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestRepository(t *testing.T) (*Repository, sqlmock.Sqlmock) {
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
	return NewRepository(database), mock
}

func TestRepositoryGetStock(t *testing.T) {
	t.Run("returns quantity", func(t *testing.T) {
		repository, mock := newTestRepository(t)
		mock.ExpectQuery("SELECT quantity FROM stock").WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(7))

		quantity, err := repository.GetStock(1)
		if err != nil || quantity != 7 {
			t.Fatalf("wanted quantity 7, got %d, %v", quantity, err)
		}
	})

	t.Run("missing stock is zero", func(t *testing.T) {
		repository, mock := newTestRepository(t)
		mock.ExpectQuery("SELECT quantity FROM stock").WithArgs(int64(2)).
			WillReturnError(sql.ErrNoRows)

		quantity, err := repository.GetStock(2)
		if err != nil || quantity != 0 {
			t.Fatalf("wanted zero quantity, got %d, %v", quantity, err)
		}
	})
}

func TestRepositoryCreateOrder(t *testing.T) {
	repository, mock := newTestRepository(t)
	createdAt := time.Date(2026, time.July, 16, 12, 0, 0, 0, time.UTC)
	items := []OrderItem{{ProductID: 1, Quantity: 2}, {ProductID: 2, Quantity: 3}}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO orders").WithArgs(StatusCreated).WillReturnRows(
		sqlmock.NewRows([]string{"id", "status", "created_at", "shipped_at"}).
			AddRow(10, StatusCreated, createdAt, nil),
	)
	for _, item := range items {
		mock.ExpectExec("INSERT INTO order_items").WithArgs(int64(10), item.ProductID, item.Quantity).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()

	created, err := repository.Create(items)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if created.ID != 10 || created.Status != StatusCreated || len(created.Items) != 2 {
		t.Fatalf("unexpected order: %#v", created)
	}
}

func TestRepositoryGetOrder(t *testing.T) {
	repository, mock := newTestRepository(t)
	createdAt := time.Date(2026, time.July, 16, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery("SELECT id,status,created_at,shipped_at FROM orders").WithArgs(int64(10)).WillReturnRows(
		sqlmock.NewRows([]string{"id", "status", "created_at", "shipped_at"}).
			AddRow(10, StatusCreated, createdAt, nil),
	)
	mock.ExpectQuery("SELECT product_id, quantity FROM order_items").WithArgs(int64(10)).WillReturnRows(
		sqlmock.NewRows([]string{"product_id", "quantity"}).
			AddRow(1, 2).
			AddRow(2, 3),
	)

	got, err := repository.Get(10)
	if err != nil {
		t.Fatalf("get order: %v", err)
	}
	if got.ID != 10 || len(got.Items) != 2 || got.Items[1].ProductID != 2 {
		t.Fatalf("unexpected order: %#v", got)
	}
}

func TestRepositoryShipOrder(t *testing.T) {
	repository, mock := newTestRepository(t)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT status FROM orders").WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(StatusCreated))
	mock.ExpectQuery("SELECT product_id, quantity FROM order_items").WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "quantity"}).AddRow(1, 2))
	mock.ExpectQuery("SELECT quantity FROM stock").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(5))
	mock.ExpectExec("UPDATE stock").WithArgs(2, int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE orders").WithArgs(StatusShipped, int64(10)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO outbound_operations").WithArgs(int64(10)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repository.Ship(10); err != nil {
		t.Fatalf("ship order: %v", err)
	}
}

func TestRepositoryRejectsAlreadyShippedOrder(t *testing.T) {
	repository, mock := newTestRepository(t)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT status FROM orders").WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(StatusShipped))
	mock.ExpectRollback()

	err := repository.Ship(10)
	if !errors.Is(err, ErrAlreadyShipped) {
		t.Fatalf("wanted already shipped error, got %v", err)
	}
}

func TestRepositoryRollsBackWhenStockIsInsufficient(t *testing.T) {
	tests := []struct {
		name      string
		stockRows *sqlmock.Rows
	}{
		{name: "stock row missing", stockRows: sqlmock.NewRows([]string{"quantity"})},
		{name: "quantity too low", stockRows: sqlmock.NewRows([]string{"quantity"}).AddRow(1)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository, mock := newTestRepository(t)
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT status FROM orders").WithArgs(int64(10)).
				WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(StatusCreated))
			mock.ExpectQuery("SELECT product_id, quantity FROM order_items").WithArgs(int64(10)).
				WillReturnRows(sqlmock.NewRows([]string{"product_id", "quantity"}).AddRow(1, 2))
			mock.ExpectQuery("SELECT quantity FROM stock").WithArgs(int64(1)).WillReturnRows(test.stockRows)
			mock.ExpectRollback()

			err := repository.Ship(10)
			if !errors.Is(err, ErrInsufficientStock) {
				t.Fatalf("wanted insufficient stock error, got %v", err)
			}
		})
	}
}
