package main

import (
	"log"
	"net/http"

	"mini-wms/internal/config"
	"mini-wms/internal/db"
	"mini-wms/internal/order"
	"mini-wms/internal/product"
	"mini-wms/internal/stock"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	productRepo := product.NewRepository(database)
	stockRepo := stock.NewRepository(database)
	orderRepo := order.NewRepository(database)

	productHandler := product.NewHandler(productRepo)
	stockHandler := stock.NewHandler(stock.NewService(stockRepo))
	orderHandler := order.NewHandler(order.NewService(orderRepo))

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.POST("/products", productHandler.Create)
	r.GET("/products", productHandler.List)
	r.POST("/inbounds", stockHandler.Inbound)
	r.GET("/stock", stockHandler.Report)
	r.POST("/orders", orderHandler.Create)
	r.GET("/orders/:id", orderHandler.Get)
	r.POST("/orders/:id/ship", orderHandler.Ship)

	log.Println("mini wms running on port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
