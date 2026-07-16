package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func orderRequest(t *testing.T, method, path, body string, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	routePath := path
	if method == http.MethodGet {
		routePath = "/orders/:id"
	}
	if method == http.MethodPost && strings.HasSuffix(path, "/ship") {
		routePath = "/orders/:id/ship"
	}
	router.Handle(method, routePath, handler)
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repo       *fakeOrderRepo
		wantStatus int
	}{
		{name: "bad JSON", body: "{", repo: newFakeOrderRepo(), wantStatus: http.StatusBadRequest},
		{name: "business error", body: `{ "items": [] }`, repo: newFakeOrderRepo(), wantStatus: http.StatusBadRequest},
		{name: "repository error", body: `{ "items": [{"product_id": 1, "quantity": 1}] }`, repo: func() *fakeOrderRepo {
			repo := newFakeOrderRepo()
			repo.stock[1] = 1
			repo.createErr = errors.New("database failed")
			return repo
		}(), wantStatus: http.StatusInternalServerError},
		{name: "created", body: `{ "items": [{"product_id": 1, "quantity": 1}] }`, repo: func() *fakeOrderRepo {
			repo := newFakeOrderRepo()
			repo.stock[1] = 1
			return repo
		}(), wantStatus: http.StatusCreated},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := NewHandler(NewService(test.repo))
			response := orderRequest(t, http.MethodPost, "/orders", test.body, handler.Create)
			if response.Code != test.wantStatus {
				t.Fatalf("wanted status %d, got %d: %s", test.wantStatus, response.Code, response.Body.String())
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	repo := newFakeOrderRepo()
	repo.orders[2] = Order{ID: 2, Status: StatusCreated}
	handler := NewHandler(NewService(repo))

	response := orderRequest(t, http.MethodGet, "/orders/:id", "", handler.Get)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("wanted bad request for invalid id, got %d", response.Code)
	}

	response = orderRequest(t, http.MethodGet, "/orders/2", "", handler.Get)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"id":2`) {
		t.Fatalf("unexpected success response: %d %s", response.Code, response.Body.String())
	}

	repo.getErr = errors.New("not found")
	response = orderRequest(t, http.MethodGet, "/orders/3", "", handler.Get)
	if response.Code != http.StatusNotFound {
		t.Fatalf("wanted not found, got %d", response.Code)
	}
}

func TestShipHandler(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		err        error
		wantStatus int
	}{
		{name: "invalid id", path: "/orders/nope/ship", wantStatus: http.StatusBadRequest},
		{name: "already shipped", path: "/orders/1/ship", err: ErrAlreadyShipped, wantStatus: http.StatusBadRequest},
		{name: "stock too low", path: "/orders/1/ship", err: ErrInsufficientStock, wantStatus: http.StatusBadRequest},
		{name: "repository error", path: "/orders/1/ship", err: errors.New("database failed"), wantStatus: http.StatusInternalServerError},
		{name: "shipped", path: "/orders/1/ship", wantStatus: http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newFakeOrderRepo()
			repo.shipErr = test.err
			handler := NewHandler(NewService(repo))
			response := orderRequest(t, http.MethodPost, test.path, "", handler.Ship)
			if response.Code != test.wantStatus {
				t.Fatalf("wanted status %d, got %d: %s", test.wantStatus, response.Code, response.Body.String())
			}
		})
	}
}

func TestParseID(t *testing.T) {
	for _, value := range []string{"", "0", "-1", "abc"} {
		if _, err := parseID(value); err == nil {
			t.Fatalf("wanted %q to be rejected", value)
		}
	}

	id, err := parseID("42")
	if err != nil || id != 42 {
		t.Fatalf("wanted id 42, got %d, %v", id, err)
	}
}
