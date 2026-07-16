package product

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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return New(db), mock
}

func TestCreateProductBadJSON(t *testing.T) {
	productDB, _ := newTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/products", productDB.Create)

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusBadRequest)
	}
}

func TestCreateProductSuccess(t *testing.T) {
	productDB, mock := newTestDB(t)

	createdAt := time.Now()
	mock.ExpectQuery("INSERT INTO products").
		WithArgs("CHAIR-1", "Chair").
		WillReturnRows(sqlmock.NewRows([]string{"id", "article_code", "name", "created_at"}).
			AddRow(1, "CHAIR-1", "Chair", createdAt))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/products", productDB.Create)

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"articleCode":"CHAIR-1","name":"Chair"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusCreated)
	}
}

func TestListProducts(t *testing.T) {
	productDB, mock := newTestDB(t)

	mock.ExpectQuery("SELECT id, article_code").
		WillReturnRows(sqlmock.NewRows([]string{"id", "article_code", "name", "created_at"}).
			AddRow(1, "CHAIR-1", "Chair", time.Now()))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/products", productDB.List)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusOK)
	}
}
