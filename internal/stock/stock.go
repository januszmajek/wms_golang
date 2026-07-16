package stock

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Models
type ReportItem struct {
	ProductID   int64  `json:"productId"`
	ArticleCode string `json:"articleCode"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
}

type InboundRequest struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type InboundResponse struct {
	ProductID     int64 `json:"productId"`
	QuantityAdded int   `json:"quantityAdded"`
}

// Database stuff
type DB struct{ db *sql.DB }

func New(database *sql.DB) *DB { return &DB{db: database} }

func (d *DB) ProductExists(productID int64) (bool, error) {
	var exists bool
	err := d.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)`, productID).Scan(&exists)
	return exists, err
}

func (d *DB) AddInbound(productID int64, quantity int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO stock (product_id, quantity) VALUES ($1,$2)
		ON CONFLICT (product_id) DO UPDATE SET quantity = stock.quantity + EXCLUDED.quantity, updated_at = NOW()`, productID, quantity)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO inbound_operations (product_id, quantity) VALUES ($1,$2)`, productID, quantity)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *DB) Report() ([]ReportItem, error) {
	rows, err := d.db.Query(`SELECT p.id, p.article_code, p.name, COALESCE(s.quantity,0) FROM products p LEFT JOIN stock s ON s.product_id=p.id ORDER BY p.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ReportItem{}
	for rows.Next() {
		var item ReportItem
		if err := rows.Scan(&item.ProductID, &item.ArticleCode, &item.Name, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// HTTP handlers
func (d *DB) Inbound(c *gin.Context) {
	var req InboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be bigger than 0"})
		return
	}

	exists, err := d.ProductExists(req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		return
	}

	if err := d.AddInbound(req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := InboundResponse{ProductID: req.ProductID, QuantityAdded: req.Quantity}
	c.JSON(http.StatusCreated, resp)
}

func (d *DB) ReportHTTP(c *gin.Context) {
	items, err := d.Report()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
