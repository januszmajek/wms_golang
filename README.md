# Mini WMS

Simple Go REST API for small warehouse. Product mgmt, inbound, stock, orders, outbound shipment, stock report. No enterprise bloat. Separated HTTP, business logic, DB.

---

## Features

- Product CRUD (create/list)
- Inbound stock receiving
- Stock stored in PostgreSQL
- Stock report showing all product inventory
- Order creation with stock validation (rejects if insufficient)
- Outbound shipment (decreases stock; prevents double shipping)
- Inbound/outbound history logging
- Unit tests (>50% coverage)

---

## Tech Stack

- **Go** (language)
- **Gin-Gonic** (web framework, middleware)
- **PostgreSQL** (persistent DB) via docker on wsl on windows and docker on Omarchy Arch Linux Distro
- **lib/pgx** (driver)
- **goose** (migrations)
- **JSON Schema** (model specification)
- **https://quicktype.io/** (golang code generation from JSON)
- **OpenAPI** (API specification)
- **Zerolog** (logger)

---

## Business Assumptions

- Backend-only REST API (no frontend, test with curl).
- Products: API-created or SQL seed.
- Inbound increases stock.
- Order creation: succeeds only if all requested items in stock. Does NOT decrease stock.
- Outbound shipment: separate endpoint, decreases stock, changes status to SHIPPED. Can't ship twice.
- No partial fulfillment.
- No Auth/Authz.
- Quantities are integers; stock cannot be negative.

---

## Main Use Cases

1. **Product Mgmt**: Create/list products (e.g. SKU `GLASS-001`).
2. **Inbound**: Add stock (increases product quantity).
3. **Stock Report**: Get all current inventory.
4. **Order Creation**: Creates `CREATED` order if stock is sufficient. No stock decrease.
5. **Shipment**: Changes status to `SHIPPED`, decreases stock, logs outbound op. Runs in DB transaction.

---

## Order Statuses

- `CREATED`: Order saved, stock reserved but not decreased.
- `SHIPPED`: Order shipped, stock decreased.
- `CANCELLED`
- `DELIVERED`: Order is successfully delivered to the customer

---

## Database Model

```sql
CREATE TABLE products (id BIGSERIAL PRIMARY KEY, sku TEXT NOT NULL UNIQUE, name TEXT NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT NOW());
CREATE TABLE stock (product_id BIGINT PRIMARY KEY REFERENCES products(id), quantity INTEGER NOT NULL CHECK (quantity >= 0), updated_at TIMESTAMP NOT NULL DEFAULT NOW());
CREATE TABLE inbound_operations (id BIGSERIAL PRIMARY KEY, product_id BIGINT NOT NULL REFERENCES products(id), quantity INTEGER NOT NULL CHECK (quantity > 0), created_at TIMESTAMP NOT NULL DEFAULT NOW());
CREATE TABLE orders (id BIGSERIAL PRIMARY KEY, status TEXT NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT NOW(), shipped_at TIMESTAMP NULL);
CREATE TABLE order_items (id BIGSERIAL PRIMARY KEY, order_id BIGINT NOT NULL REFERENCES orders(id), product_id BIGINT NOT NULL REFERENCES products(id), quantity INTEGER NOT NULL CHECK (quantity > 0));
CREATE TABLE outbound_operations (id BIGSERIAL PRIMARY KEY, order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id), created_at TIMESTAMP NOT NULL DEFAULT NOW());
```

---

## API Endpoints

- **Create Product**: `POST /products` `{"sku":"GLASS-001","name":"Glass"}` -> `{"id":1,...}`
- **List Products**: `GET /products` -> `[{"id":1,...}]`
- **Inbound**: `POST /inbounds` `{"product_id":1,"quantity":10}` -> `{"product_id":1,"quantity_added":10}`
- **Stock Report**: `GET /stock` -> `[{"product_id":1,"sku":"GLASS-001","name":"Glass","quantity":10}]`
- **Create Order**: `POST /orders` `{"items":[{"product_id":1,"quantity":4}]}` -> `{"id":1,"status":"CREATED",...}` (400 if insufficient)
- **Get Order**: `GET /orders/:id` -> `{"id":1,...}`
- **Change Order Status**: `POST /orders/:id/<status>` -> `{"order_id":1,"status":"SHIPPED"}` (errors: already shipped, stock changed etc.)

---

## Business Rules

- **Order creation stock validation**: Sum duplicate products in request before checking stock. E.g. request `[{id:1, qty:6}, {id:1, qty:6}]` requires stock >= 12.
- **Shipment transaction**: Re-check stock during shipment (may change after creation). Wrap in DB txn:
  1. Check status. 2. Check stock. 3. Decrease stock. 4. Update order to SHIPPED. 5. Record outbound operation.

---

## Out of Scope

Frontend, user accounts, auth, categories, reservations, partial shipments, order cancellation, pagination, advanced filtering, multiple warehouses, event queues, external integrations.

---

## Assignment Coverage

Covers inbound stock, orders validation, outbound shipment, inventory report, PostgreSQL persistence, and unit tests with >50% coverage.

---
