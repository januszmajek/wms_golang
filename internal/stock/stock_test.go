package stock

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestInboundHandlerBadJSON(t *testing.T) {
	stockDB, _ := newTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/inbounds", stockDB.Inbound)

	req := httptest.NewRequest(http.MethodPost, "/inbounds", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusBadRequest)
	}
}

func TestInboundHandlerSuccess(t *testing.T) {
	stockDB, mock := newTestDB(t)

	mock.ExpectQuery("SELECT EXISTS").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO stock").WithArgs(int64(1), 5).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO inbound_operations").WithArgs(int64(1), 5).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/inbounds", stockDB.Inbound)

	req := httptest.NewRequest(http.MethodPost, "/inbounds", strings.NewReader(`{"productId":1,"quantity":5}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusCreated)
	}
}

func TestReportHandler(t *testing.T) {
	stockDB, mock := newTestDB(t)

	mock.ExpectQuery("SELECT p.id").
		WillReturnRows(sqlmock.NewRows([]string{"id", "article_code", "name", "quantity"}).
			AddRow(1, "TEST-1", "Test", 10))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/stock", stockDB.ReportHTTP)

	req := httptest.NewRequest(http.MethodGet, "/stock", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.Code, http.StatusOK)
	}
}
