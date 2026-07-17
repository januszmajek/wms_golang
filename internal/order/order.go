package order

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const StatusCreated = "CREATED"
const StatusShipped = "SHIPPED"

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrAlreadyShipped    = errors.New("order already shipped")
)

type ItemRequest struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type CreateOrderRequest struct {
	Items       []ItemRequest `json:"items"`
	Description string        `json:"description,omitempty"`
}

type OrderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type Order struct {
	ID          int64       `json:"id"`
	Status      string      `json:"status"`
	Description string      `json:"description,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	ShippedAt   *time.Time  `json:"shippedAt,omitempty"`
	Items       []OrderItem `json:"items"`
}

type ShipResponse struct {
	OrderID int64  `json:"orderId"`
	Status  string `json:"status"`
}

type DB struct{ db *sql.DB }

func New(database *sql.DB) *DB { return &DB{db: database} }

func (d *DB) GetStock(productID int64) (int, error) {
	var qty int
	err := d.db.QueryRow(`SELECT quantity FROM stock WHERE product_id=$1`, productID).Scan(&qty)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return qty, err
}

func (d *DB) CreateOrder(items []OrderItem, description string) (Order, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()
	var o Order
	err = tx.QueryRow(`INSERT INTO orders (status, description) VALUES ($1, $2) RETURNING id,status,description,created_at,shipped_at`, StatusCreated, description).Scan(&o.ID, &o.Status, &o.Description, &o.CreatedAt, &o.ShippedAt)
	if err != nil {
		return Order{}, err
	}
	for _, item := range items {
		_, err = tx.Exec(`INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1,$2,$3)`, o.ID, item.ProductID, item.Quantity)
		if err != nil {
			return Order{}, err
		}
	}
	if err = tx.Commit(); err != nil {
		return Order{}, err
	}
	o.Items = items
	return o, nil
}

func (d *DB) GetOrder(id int64) (Order, error) {
	var o Order
	err := d.db.QueryRow(`SELECT id,status,description,created_at,shipped_at FROM orders WHERE id=$1`, id).Scan(&o.ID, &o.Status, &o.Description, &o.CreatedAt, &o.ShippedAt)
	if err != nil {
		return Order{}, err
	}
	rows, err := d.db.Query(`SELECT product_id, quantity FROM order_items WHERE order_id=$1 ORDER BY id`, id)
	if err != nil {
		return Order{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return Order{}, err
		}
		o.Items = append(o.Items, item)
	}
	return o, rows.Err()
}

func (d *DB) ShipOrder(id int64) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var status string
	err = tx.QueryRow(`SELECT status FROM orders WHERE id=$1 FOR UPDATE`, id).Scan(&status)
	if err != nil {
		return err
	}
	if status == StatusShipped {
		return ErrAlreadyShipped
	}
	rows, err := tx.Query(`SELECT product_id, quantity FROM order_items WHERE order_id=$1`, id)
	if err != nil {
		return err
	}
	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			err := rows.Close()
			if err != nil {
				return err
			}
		}
		items = append(items, item)
	}
	rows.Close()
	for _, item := range items {
		var stock int
		err = tx.QueryRow(`SELECT quantity FROM stock WHERE product_id=$1 FOR UPDATE`, item.ProductID).Scan(&stock)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInsufficientStock
		}
		if err != nil {
			return err
		}
		if stock < item.Quantity {
			return ErrInsufficientStock
		}
		_, err = tx.Exec(`UPDATE stock SET quantity=quantity-$1, updated_at=NOW() WHERE product_id=$2`, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(`UPDATE orders SET status=$1, shipped_at=NOW() WHERE id=$2`, StatusShipped, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO outbound_operations (order_id) VALUES ($1)`, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *DB) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order needs at least one item"})
		return
	}

	merged := map[int64]int{}
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "item quantity must be bigger than 0"})
			return
		}
		merged[item.ProductID] += item.Quantity
	}

	var items []OrderItem
	for productID, qty := range merged {
		stock, err := d.GetStock(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if stock < qty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock"})
			return
		}
		items = append(items, OrderItem{ProductID: productID, Quantity: qty})
	}

	order, err := d.CreateOrder(items, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (d *DB) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}
	o, err := d.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (d *DB) Ship(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	err = d.ShipOrder(id)
	if errors.Is(err, ErrAlreadyShipped) || errors.Is(err, ErrInsufficientStock) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := ShipResponse{OrderID: id, Status: StatusShipped}
	c.JSON(http.StatusOK, resp)
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
