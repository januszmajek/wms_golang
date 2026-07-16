package stock

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func stockRequest(t *testing.T, method, path, body string, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(method, path, handler)
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestInboundHandler(t *testing.T) {
	t.Run("rejects bad JSON", func(t *testing.T) {
		service, _ := newTestService(t)
		response := stockRequest(t, http.MethodPost, "/inbounds", "{", NewHandler(service).Inbound)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("wanted bad request, got %d", response.Code)
		}
	})

	t.Run("returns business error", func(t *testing.T) {
		service, _ := newTestService(t)
		response := stockRequest(t, http.MethodPost, "/inbounds", `{"product_id":1,"quantity":0}`, NewHandler(service).Inbound)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("wanted bad request, got %d", response.Code)
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		service, mock := newTestService(t)
		mock.ExpectQuery("SELECT EXISTS").WithArgs(int64(1)).WillReturnError(errors.New("database unavailable"))
		response := stockRequest(t, http.MethodPost, "/inbounds", `{"product_id":1,"quantity":2}`, NewHandler(service).Inbound)
		if response.Code != http.StatusInternalServerError {
			t.Fatalf("wanted server error, got %d", response.Code)
		}
	})

	t.Run("receives stock", func(t *testing.T) {
		service, mock := newTestService(t)
		mock.ExpectQuery("SELECT EXISTS").WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO stock").WithArgs(int64(1), 2).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO inbound_operations").WithArgs(int64(1), 2).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		response := stockRequest(t, http.MethodPost, "/inbounds", `{"product_id":1,"quantity":2}`, NewHandler(service).Inbound)
		if response.Code != http.StatusCreated || !strings.Contains(response.Body.String(), `"quantity_added":2`) {
			t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
		}
	})
}

func TestReportHandler(t *testing.T) {
	t.Run("returns report", func(t *testing.T) {
		service, mock := newTestService(t)
		mock.ExpectQuery("SELECT p.id").WillReturnRows(
			sqlmock.NewRows([]string{"id", "article_code", "name", "quantity"}).AddRow(1, "CHAIR-1", "Chair", 3),
		)
		response := stockRequest(t, http.MethodGet, "/stock", "", NewHandler(service).Report)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"quantity":3`) {
			t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		service, mock := newTestService(t)
		mock.ExpectQuery("SELECT p.id").WillReturnError(errors.New("database unavailable"))
		response := stockRequest(t, http.MethodGet, "/stock", "", NewHandler(service).Report)
		if response.Code != http.StatusInternalServerError {
			t.Fatalf("wanted server error, got %d", response.Code)
		}
	})
}
