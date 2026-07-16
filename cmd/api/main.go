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

	productDB := product.New(database)
	stockDB := stock.New(database)
	orderDB := order.New(database)

	r := newRouter(productDB, stockDB, orderDB)

	log.Println("mini wms running on port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}

func newRouter(productDB *product.DB, stockDB *stock.DB, orderDB *order.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/stock", stockDB.ReportHTTP)
	r.GET("/products", productDB.List)
	r.GET("/orders/:id", orderDB.Get)
	r.POST("/products", productDB.Create)
	r.POST("/inbounds", stockDB.Inbound)
	r.POST("/orders", orderDB.Create)
	r.POST("/orders/:id/ship", orderDB.Ship)

	return r
}
