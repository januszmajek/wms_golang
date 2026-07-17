-- +goose Up
ALTER TABLE orders ADD COLUMN IF NOT EXISTS description TEXT;

-- +goose Down
ALTER TABLE orders DROP COLUMN IF EXISTS description;
