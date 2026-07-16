package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"mini-wms/internal/order"
	"mini-wms/internal/product"
	"mini-wms/internal/stock"

	"github.com/gin-gonic/gin"
)

func TestRouterHealthAndRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := newRouter(
		product.NewHandler(nil),
		stock.NewHandler(nil),
		order.NewHandler(nil),
	)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Body.String() != `{"status":"ok"}` {
		t.Fatalf("unexpected health response: %d %s", response.Code, response.Body.String())
	}

	wantRoutes := map[string]bool{
		"GET /health":           false,
		"GET /stock":            false,
		"GET /products":         false,
		"GET /orders/:id":       false,
		"POST /products":        false,
		"POST /inbounds":        false,
		"POST /orders":          false,
		"POST /orders/:id/ship": false,
	}
	for _, route := range router.Routes() {
		key := route.Method + " " + route.Path
		if _, exists := wantRoutes[key]; exists {
			wantRoutes[key] = true
		}
	}
	for route, registered := range wantRoutes {
		if !registered {
			t.Errorf("route %s is not registered", route)
		}
	}
}
