-- +goose Up
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  product_code VARCHAR(50) UNIQUE NOT NULL,
  name VARCHAR(75) NOT NULL
);

-- +goose Down
DROP TABLE products;