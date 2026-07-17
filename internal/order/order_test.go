package order

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func newTestDB(t *testing.T) (*DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("failed to create sql mock: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return New(db), mock
}

func TestCreateHandlerBadJSON(t *testing.T) {
	orderDB, _ := newTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/orders", orderDB.Create)

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusBadRequest)
	}
}

func TestCreateHandlerSuccess(t *testing.T) {
	orderDB, mock := newTestDB(t)

	mock.ExpectQuery("SELECT quantity FROM stock").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(10))

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO orders").WithArgs(StatusCreated, "").
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "description", "created_at", "shipped_at"}).
			AddRow(1, StatusCreated, "", time.Now(), nil))
	mock.ExpectExec("INSERT INTO order_items").WithArgs(1, int64(1), 5).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/orders", orderDB.Create)

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"items":[{"productId":1,"quantity":5}]}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusCreated)
	}
}

func TestGetHandler(t *testing.T) {
	orderDB, mock := newTestDB(t)

	mock.ExpectQuery("SELECT id,status,description,created_at,shipped_at FROM orders").WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "description", "created_at", "shipped_at"}).
			AddRow(5, StatusCreated, "", time.Now(), nil))
	mock.ExpectQuery("SELECT product_id, quantity FROM order_items").WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "quantity"}).
			AddRow(int64(1), 3))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/orders/:id", orderDB.Get)

	req := httptest.NewRequest(http.MethodGet, "/orders/5", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusOK)
	}
}

func TestShipHandlerSuccess(t *testing.T) {
	orderDB, mock := newTestDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT status FROM orders").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(StatusCreated))
	mock.ExpectQuery("SELECT product_id, quantity FROM order_items").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "quantity"}).
			AddRow(int64(1), 5))
	mock.ExpectQuery("SELECT quantity FROM stock").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(10))
	mock.ExpectExec("UPDATE stock").WithArgs(5, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE orders").WithArgs(StatusShipped, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO outbound_operations").WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/orders/:id/ship", orderDB.Ship)

	req := httptest.NewRequest(http.MethodPost, "/orders/1/ship", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusOK)
	}
}

func TestParseID(t *testing.T) {
	_, err := parseID("abc")
	if err == nil {
		t.Error("expected error for invalid ID")
	}

	id, err := parseID("42")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Errorf("got ID %d, want 42", id)
	}
}
