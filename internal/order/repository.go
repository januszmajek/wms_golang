package order

import (
	"database/sql"
	"errors"
)

type Repository struct{ DB *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{DB: db} }
func (r *Repository) GetStock(productID int64) (int, error) {
	var qty int
	err := r.DB.QueryRow(`SELECT quantity FROM stock WHERE product_id=$1`, productID).Scan(&qty)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return qty, err
}
func (r *Repository) Create(items []OrderItem) (Order, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()
	var o Order
	err = tx.QueryRow(`INSERT INTO orders (status) VALUES ($1) RETURNING id,status,created_at,shipped_at`, StatusCreated).Scan(&o.ID, &o.Status, &o.CreatedAt, &o.ShippedAt)
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
func (r *Repository) Get(id int64) (Order, error) {
	var o Order
	err := r.DB.QueryRow(`SELECT id,status,created_at,shipped_at FROM orders WHERE id=$1`, id).Scan(&o.ID, &o.Status, &o.CreatedAt, &o.ShippedAt)
	if err != nil {
		return Order{}, err
	}
	rows, err := r.DB.Query(`SELECT product_id, quantity FROM order_items WHERE order_id=$1 ORDER BY id`, id)
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
func (r *Repository) Ship(id int64) error {
	tx, err := r.DB.Begin()
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
	items := []OrderItem{}
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			rows.Close()
			return err
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
