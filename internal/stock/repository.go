package stock

import (
	"database/sql"
)

type Repository struct{ DB *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{DB: db} }

func (r *Repository) ProductExists(productID int64) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)`, productID).Scan(&exists)
	return exists, err
}

func (r *Repository) AddInbound(productID int64, quantity int) error {
	tx, err := r.DB.Begin()
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

func (r *Repository) Report() ([]ReportItem, error) {
	rows, err := r.DB.Query(`SELECT p.id, p.article_code, p.name, COALESCE(s.quantity,0) FROM products p LEFT JOIN stock s ON s.product_id=p.id ORDER BY p.id`)
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
