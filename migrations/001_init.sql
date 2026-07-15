-- +goose Up
CREATE TABLE products (id BIGSERIAL PRIMARY KEY, article_code TEXT NOT NULL UNIQUE, name TEXT NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT NOW());
CREATE TABLE stock (product_id BIGINT PRIMARY KEY REFERENCES products(id), quantity INTEGER NOT NULL CHECK (quantity >= 0), updated_at TIMESTAMP NOT NULL DEFAULT NOW());
CREATE TABLE inbound_operations (id BIGSERIAL PRIMARY KEY, product_id BIGINT NOT NULL REFERENCES products(id), quantity INTEGER NOT NULL CHECK (quantity > 0), created_at TIMESTAMP NOT NULL DEFAULT NOW());
CREATE TABLE orders (id BIGSERIAL PRIMARY KEY, status TEXT NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT NOW(), shipped_at TIMESTAMP NULL);
CREATE TABLE order_items (id BIGSERIAL PRIMARY KEY, order_id BIGINT NOT NULL REFERENCES orders(id), product_id BIGINT NOT NULL REFERENCES products(id), quantity INTEGER NOT NULL CHECK (quantity > 0));
CREATE TABLE outbound_operations (id BIGSERIAL PRIMARY KEY, order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id), created_at TIMESTAMP NOT NULL DEFAULT NOW());

-- +goose Down
DROP TABLE outbound_operations;
DROP TABLE order_items;
DROP TABLE orders;
DROP TABLE inbound_operations;
DROP TABLE stock;
DROP TABLE products;
