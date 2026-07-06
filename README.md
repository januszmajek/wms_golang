# Mini WMS

Mini WMS is a simple monolithic REST API written in Go for managing a small warehouse. The application supports basic warehouse operations: product management, inbound stock receiving, order creation with stock validation, outbound shipment, and current inventory reporting.

The project was designed as a small learning-oriented backend assignment. It intentionally avoids unnecessary enterprise complexity while still keeping a clean separation between HTTP handlers, business logic, and database access.

---

## Features

- Create and list products.
- Receive products into the warehouse through inbound operations.
- Store current stock levels in PostgreSQL.
- Generate a stock report with current inventory for all products.
- Create customer orders for one or more products.
- Validate whether an order can be fulfilled based on current stock.
- Reject order creation when requested quantity is greater than available stock.
- Ship accepted orders through a separate outbound operation.
- Decrease stock only when an order is shipped.
- Prevent shipping the same order more than once.
- Keep basic operation history for inbound and outbound actions.
- Include unit tests for business logic with code coverage above 50%.

---

## Tech Stack

- **Go** — main programming language.
- **Gin** — HTTP web framework used to build the REST API.
- **PostgreSQL** — relational database used for persistent storage.
- **Podman Desktop** — used to run PostgreSQL locally.
- **database/sql** — Go standard package used for database access.
- **lib/pq** — PostgreSQL driver for Go.
- **goose** — database migration tool.
- **testing** — Go standard testing package for unit tests.

---

## Architecture

The application follows a simple monolithic architecture:

```text
HTTP handlers -> services -> repositories -> PostgreSQL
```

### Handlers

Handlers are responsible for HTTP-specific concerns:

- parsing JSON requests,
- validating basic input format,
- calling the correct service method,
- returning JSON responses and HTTP status codes.

### Services

Services contain business logic, for example:

- checking if inbound quantity is valid,
- validating whether an order can be created,
- checking current stock,
- preventing negative stock,
- shipping orders,
- preventing double shipment.

This layer is the main target for unit tests.

### Repositories

Repositories are responsible for database operations:

- inserting records,
- selecting records,
- updating stock quantities,
- creating orders and order items,
- using transactions where needed.

---

## Business Assumptions

- The application is a backend-only REST API. It does not include a frontend.
- The API can be tested with Postman or curl.
- Products may be created through the API or prepared using SQL seed data.
- Inbound operations increase product stock.
- An order can be created only if every requested product has enough available stock.
- Stock is not decreased when an order is created.
- Stock is decreased only during outbound shipment.
- An order must be shipped using a separate endpoint.
- An already shipped order cannot be shipped again.
- Partial order fulfillment is not supported.
- Authentication and authorization are not included.
- Product quantities are represented as integers.
- Stock quantity must never become negative.

---

## Main Use Cases

### 1. Product Management

The user can create products that will later be used in inbound operations and orders.

Example products:

- `GLASS-001` — Glass
- `PLATE-001` — Plate
- `MUG-001` — Mug

### 2. Inbound

Inbound represents receiving products into the warehouse.

Example:

```text
Receive 10 glasses into stock.
```

After this operation, the stock quantity for the selected product increases by 10.

### 3. Inventory / Stock Report

The user can request the current warehouse inventory.

The stock report returns all products with their current quantities.

### 4. Order Creation

The user can create an order containing one or more products and requested quantities.

Before the order is saved, the application checks whether all requested quantities are available.

If stock is sufficient, the order is created with status `CREATED`.

If stock is insufficient, the API returns an error and the order is not created.

### 5. Outbound / Shipment

Outbound represents shipping an existing order.

When an order is shipped:

- the order status changes from `CREATED` to `SHIPPED`,
- stock is decreased by the quantities from the order,
- an outbound operation is recorded.

The shipment operation should be executed in a database transaction.

---

## Order Statuses

The application uses a minimal set of order statuses:

```text
CREATED
SHIPPED
```

### CREATED

The order was successfully created and can be shipped.

### SHIPPED

The order was shipped and stock was decreased.

---

## Database Model

### `products`

Stores product catalog data.

```sql
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    sku TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### `stock`

Stores current stock quantity for each product.

```sql
CREATE TABLE stock (
    product_id BIGINT PRIMARY KEY REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### `inbound_operations`

Stores history of received stock.

```sql
CREATE TABLE inbound_operations (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### `orders`

Stores order headers.

```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    shipped_at TIMESTAMP NULL
);
```

### `order_items`

Stores products and quantities assigned to each order.

```sql
CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0)
);
```

### `outbound_operations`

Stores shipment history.

```sql
CREATE TABLE outbound_operations (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## API Endpoints

### Health Check

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

---

### Create Product

```http
POST /products
```

Request:

```json
{
  "sku": "GLASS-001",
  "name": "Glass"
}
```

Response:

```json
{
  "id": 1,
  "sku": "GLASS-001",
  "name": "Glass"
}
```

---

### List Products

```http
GET /products
```

Response:

```json
[
  {
    "id": 1,
    "sku": "GLASS-001",
    "name": "Glass"
  }
]
```

---

### Create Inbound Operation

```http
POST /inbounds
```

Request:

```json
{
  "product_id": 1,
  "quantity": 10
}
```

Response:

```json
{
  "product_id": 1,
  "quantity_added": 10
}
```

Result:

```text
Stock for product_id=1 is increased by 10.
```

---

### Get Stock Report

```http
GET /stock
```

Response:

```json
[
  {
    "product_id": 1,
    "sku": "GLASS-001",
    "name": "Glass",
    "quantity": 10
  }
]
```

---

### Create Order

```http
POST /orders
```

Request:

```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 4
    }
  ]
}
```

Response when stock is sufficient:

```json
{
  "id": 1,
  "status": "CREATED",
  "items": [
    {
      "product_id": 1,
      "quantity": 4
    }
  ]
}
```

Response when stock is insufficient:

```json
{
  "error": "insufficient stock for product_id=1, requested=20, available=10"
}
```

Expected HTTP status:

```text
400 Bad Request
```

Important: creating an order does not decrease stock. Stock is decreased only when the order is shipped.

---

### Get Order By ID

```http
GET /orders/:id
```

Response:

```json
{
  "id": 1,
  "status": "CREATED",
  "items": [
    {
      "product_id": 1,
      "quantity": 4
    }
  ]
}
```

This endpoint is optional, but useful for debugging and presentation.

---

### Ship Order

```http
POST /orders/:id/ship
```

Response:

```json
{
  "order_id": 1,
  "status": "SHIPPED"
}
```

Result:

```text
The order is marked as shipped and stock is decreased.
```

Error when the order has already been shipped:

```json
{
  "error": "order already shipped"
}
```

Error when stock is no longer available during shipment:

```json
{
  "error": "insufficient stock during shipment"
}
```

---

## Example Demo Flow

This is the main flow that demonstrates all required business cases.

### 1. Create a product

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"sku":"GLASS-001","name":"Glass"}'
```

### 2. Receive stock

```bash
curl -X POST http://localhost:8080/inbounds \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":10}'
```

### 3. Check stock

```bash
curl http://localhost:8080/stock
```

Expected result:

```text
Product GLASS-001 has quantity 10.
```

### 4. Create an order for available quantity

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":4}]}'
```

Expected result:

```text
Order is created with status CREATED.
```

### 5. Ship the order

```bash
curl -X POST http://localhost:8080/orders/1/ship
```

Expected result:

```text
Order status changes to SHIPPED.
Stock decreases from 10 to 6.
```

### 6. Check stock again

```bash
curl http://localhost:8080/stock
```

Expected result:

```text
Product GLASS-001 has quantity 6.
```

### 7. Try to create an order with insufficient stock

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":20}]}'
```

Expected result:

```text
API returns an insufficient stock error.
```

---

## Local Development Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd mini-wms
```

### 2. Start PostgreSQL

```bash
docker compose up -d
```

### 3. Run database migrations

```bash
goose -dir migrations postgres "postgres://wms:wms@localhost:5432/mini_wms?sslmode=disable" up
```

### 4. Run the application

```bash
go run ./cmd/api
```

The API should be available at:

```text
http://localhost:8080
```

---

## Environment Variables

Example `.env.example`:

```env
APP_PORT=8080
DATABASE_URL=postgres://wms:wms@localhost:5432/mini_wms?sslmode=disable
```

If environment variables are not implemented, the application can use the same values as defaults in the configuration code.

---

## Database Migrations

The project uses `goose` for database migrations.

Create a new migration:

```bash
goose -dir migrations create migration_name sql
```

Run migrations:

```bash
goose -dir migrations postgres "postgres://wms:wms@localhost:5432/mini_wms?sslmode=disable" up
```

Rollback the last migration:

```bash
goose -dir migrations postgres "postgres://wms:wms@localhost:5432/mini_wms?sslmode=disable" down
```

---

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test ./... -cover
```

Generate an HTML coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

The required code coverage is above 50%.

---

## Recommended Unit Test Coverage

The most important tests should cover the service layer.

### Stock Service

- Inbound operation increases stock.
- Inbound operation rejects quantity less than or equal to zero.
- Inbound operation fails for a non-existing product.

### Order Service

- Order is created when stock is sufficient.
- Order creation fails when stock is insufficient.
- Order creation fails for an empty item list.
- Order creation fails when item quantity is less than or equal to zero.
- Order creation correctly handles duplicate product IDs in one request.
- Order creation works for multiple products.
- Order creation fails when at least one product has insufficient stock.
- Shipping an order decreases stock.
- Shipping an order changes its status to `SHIPPED`.
- Already shipped order cannot be shipped again.
- Shipment fails when stock is no longer sufficient at shipment time.

---

## Important Business Rules

### Stock validation during order creation

When an order is created, the application must check whether all requested products are available in the requested quantities.

If the request contains the same product multiple times, quantities should be summed before stock validation.

Example:

```json
{
  "items": [
    { "product_id": 1, "quantity": 6 },
    { "product_id": 1, "quantity": 6 }
  ]
}
```

If the available stock for product `1` is `10`, this order must be rejected because the total requested quantity is `12`.

### Stock validation during shipment

Stock should be checked again during shipment because stock may have changed between order creation and shipment.

Shipment should run in a database transaction to avoid partial updates.

The transaction should include:

1. checking order status,
2. checking current stock,
3. decreasing stock,
4. updating order status,
5. creating outbound operation record.

---

## Suggested Project Structure

```text
mini-wms/
  cmd/
    api/
      main.go

  internal/
    config/
      config.go

    db/
      db.go

    product/
      model.go
      repository.go
      handler.go

    stock/
      model.go
      repository.go
      service.go
      handler.go
      service_test.go

    order/
      model.go
      repository.go
      service.go
      handler.go
      service_test.go

  migrations/
    001_init.sql

  docker-compose.yml
  .env.example
  go.mod
  go.sum
  README.md
```

A simpler structure is also acceptable for this assignment, as long as the code remains readable and business logic is not hidden inside HTTP handlers.

---

## Scope Not Included

The following features are intentionally not implemented:

- frontend application,
- user accounts,
- authentication,
- authorization,
- product categories,
- reservations,
- partial shipments,
- order cancellation,
- pagination,
- advanced filtering,
- warehouse locations,
- multiple warehouses,
- event queues,
- external integrations.

---

## Assignment Coverage

This project covers the required assignment scope:

- **Inbound** — products can be received into warehouse stock.
- **Order** — orders can be created only when requested stock is available.
- **Outbound** — accepted orders can be shipped and stock is decreased.
- **Stock / Inventory** — current stock report can be generated.
- **Database** — data is persisted in PostgreSQL.
- **Unit tests** — business logic is covered with tests and coverage should be above 50%.

---

## Final Demo Scenario

A complete presentation can use this sequence:

```text
1. Start PostgreSQL with Docker Compose.
2. Run goose migrations.
3. Start the Go API.
4. Create product GLASS-001.
5. Add inbound quantity 10.
6. Display stock report showing quantity 10.
7. Create order for quantity 4.
8. Ship the order.
9. Display stock report showing quantity 6.
10. Try to create order for quantity 20.
11. Show insufficient stock error.
12. Run unit tests with coverage.
```

This flow demonstrates all core requirements of the mini warehouse management system.
