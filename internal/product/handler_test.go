package product

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func newProductHandler(t *testing.T) (*Handler, sqlmock.Sqlmock) {
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
	return NewHandler(NewRepository(database)), mock
}

func productRequest(t *testing.T, method, body string, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(method, "/products", handler)
	request := httptest.NewRequest(method, "/products", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestCreateProductHandler(t *testing.T) {
	t.Run("rejects bad JSON", func(t *testing.T) {
		handler, _ := newProductHandler(t)
		response := productRequest(t, http.MethodPost, "{", handler.Create)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("wanted bad request, got %d", response.Code)
		}
	})

	t.Run("rejects blank fields", func(t *testing.T) {
		handler, _ := newProductHandler(t)
		response := productRequest(t, http.MethodPost, `{"article_code":"  ","name":"Chair"}`, handler.Create)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("wanted bad request, got %d", response.Code)
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		handler, mock := newProductHandler(t)
		mock.ExpectQuery("INSERT INTO products").WithArgs("CHAIR-1", "Chair").
			WillReturnError(errors.New("duplicate article code"))
		response := productRequest(t, http.MethodPost, `{"article_code":"CHAIR-1","name":"Chair"}`, handler.Create)
		if response.Code != http.StatusInternalServerError {
			t.Fatalf("wanted server error, got %d", response.Code)
		}
	})

	t.Run("creates and trims product", func(t *testing.T) {
		handler, mock := newProductHandler(t)
		createdAt := time.Date(2026, time.July, 16, 12, 0, 0, 0, time.UTC)
		mock.ExpectQuery("INSERT INTO products").WithArgs("CHAIR-1", "Chair").WillReturnRows(
			sqlmock.NewRows([]string{"id", "article_code", "name", "created_at"}).
				AddRow(1, "CHAIR-1", "Chair", createdAt),
		)
		response := productRequest(t, http.MethodPost, `{"article_code":" CHAIR-1 ","name":" Chair "}`, handler.Create)
		if response.Code != http.StatusCreated || !strings.Contains(response.Body.String(), `"article_code":"CHAIR-1"`) {
			t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
		}
	})
}

func TestListProductsHandler(t *testing.T) {
	t.Run("lists products", func(t *testing.T) {
		handler, mock := newProductHandler(t)
		mock.ExpectQuery("SELECT id, article_code").WillReturnRows(
			sqlmock.NewRows([]string{"id", "article_code", "name", "created_at"}).
				AddRow(1, "CHAIR-1", "Chair", time.Now()).
				AddRow(2, "TABLE-1", "Table", time.Now()),
		)
		response := productRequest(t, http.MethodGet, "", handler.List)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"TABLE-1"`) {
			t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
		}
	})

	t.Run("returns query error", func(t *testing.T) {
		handler, mock := newProductHandler(t)
		mock.ExpectQuery("SELECT id, article_code").WillReturnError(errors.New("database unavailable"))
		response := productRequest(t, http.MethodGet, "", handler.List)
		if response.Code != http.StatusInternalServerError {
			t.Fatalf("wanted server error, got %d", response.Code)
		}
	})
}
