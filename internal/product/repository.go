package product

import "database/sql"

type Repository struct{ DB *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{DB: db} }

func (r *Repository) Create(req CreateProductRequest) (Product, error) {
	var p Product
	err := r.DB.QueryRow(`INSERT INTO products (sku, name) VALUES ($1, $2) RETURNING id, sku, name, created_at`, req.SKU, req.Name).
		Scan(&p.ID, &p.SKU, &p.Name, &p.CreatedAt)
	return p, err
}

func (r *Repository) List() ([]Product, error) {
	rows, err := r.DB.Query(`SELECT id, sku, name, created_at FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}
