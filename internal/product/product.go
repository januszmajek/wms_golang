package product

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Models
type Product struct {
	ID          int64     `json:"id"`
	ArticleCode string    `json:"articleCode"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateProductRequest struct {
	ArticleCode string `json:"articleCode"`
	Name        string `json:"name"`
}

// Database stuff
type DB struct{ db *sql.DB }

func New(database *sql.DB) *DB { return &DB{db: database} }

func (d *DB) CreateProduct(req CreateProductRequest) (Product, error) {
	var p Product
	err := d.db.QueryRow(`INSERT INTO products (article_code, name) VALUES ($1, $2) RETURNING id, article_code, name, created_at`, req.ArticleCode, req.Name).
		Scan(&p.ID, &p.ArticleCode, &p.Name, &p.CreatedAt)
	return p, err
}

func (d *DB) ListProducts() ([]Product, error) {
	rows, err := d.db.Query(`SELECT id, article_code, name, created_at FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.ArticleCode, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// HTTP handlers
func (d *DB) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	req.ArticleCode = strings.TrimSpace(req.ArticleCode)
	req.Name = strings.TrimSpace(req.Name)
	if req.ArticleCode == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "article_code and name are required"})
		return
	}

	p, err := d.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (d *DB) List(c *gin.Context) {
	products, err := d.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
