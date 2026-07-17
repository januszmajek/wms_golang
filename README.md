# Mini WMS

Mini WMS is a small Go REST API for managing products, warehouse stock, and outbound orders. It is a backend-only example project: there is no frontend, authentication, authorization, or external integration.

## What it does

- Creates and lists products.
- Receives inbound stock and records each inbound operation.
- Reports the current quantity for every product, including products with zero stock.
- Creates orders after checking that all requested quantities are currently in stock.
- Ships orders in a database transaction, decreases stock, records the outbound operation, and prevents a second shipment.

An order is created with status `CREATED`. Creating an order does not decrease stock. Shipping re-checks stock and changes the order to `SHIPPED`.

## Technology

- Go 1.24
- Gin 1.10 for HTTP routing and JSON handling
- PostgreSQL 16
- `database/sql` with `lib/pq`
- SQL migrations in `migrations/001_init.sql`
- Go's standard `testing` package, small fakes, and `go-sqlmock`

## Prerequisites

- Go 1.24 or newer
- Podman with a running Podman machine, or Docker
- PostgreSQL client tools are optional; they are only needed for direct database commands

## Quick start

From the repository root, start PostgreSQL with the included Compose file:

```powershell
podman machine start podman-machine-default
podman compose up -d
```

The database is exposed on `localhost:5433` and uses:

| Setting | Value |
| --- | --- |
| Database | `mini_wms` |
| User | `postgres` |
| Password | `postgres` |
| Host port | `5433` |

Apply the migration. The project does not currently run migrations automatically when the API starts:

```powershell
go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable" up
```

Start the API in a second terminal:

```powershell
go run ./cmd/api
```

The API listens on `http://localhost:8081` by default. Verify it with:

```powershell
curl http://localhost:8081/health
```

Expected response:

```json
{"status":"ok"}
```

To stop PostgreSQL:

```powershell
podman compose down
```

## Configuration

Configuration is read directly from environment variables. No `.env` file is loaded by the application.

| Variable | Default | Description |
| --- | --- | --- |
| `APP_PORT` | `8081` | HTTP listen port |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable` | PostgreSQL connection string |

Example:

```powershell
$env:APP_PORT = "8082"
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable"
go run ./cmd/api
```

## API

All request and response bodies are JSON. Validation errors normally return HTTP `400` with `{"error":"..."}`. Database and unexpected errors return HTTP `500`.

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/health` | Returns `{"status":"ok"}`. |
| `POST` | `/products` | Creates a product. Required fields: `articleCode`, `name`. Returns `201`. |
| `GET` | `/products` | Lists products ordered by ID. |
| `POST` | `/inbounds` | Adds a positive `quantity` to an existing `productId`. Returns `201`. |
| `GET` | `/stock` | Lists every product and its current quantity. |
| `POST` | `/orders` | Creates an order with a non-empty `items` array. Each item needs a positive `productId` and `quantity`. Returns `201` with status `CREATED`. |
| `GET` | `/orders/:id` | Returns an order and its items. The ID must be a positive integer; missing database IDs return `404`. |
| `POST` | `/orders/:id/ship` | Re-checks stock, decreases it, and ships the order. The ID must be a positive integer. Returns status `SHIPPED`. |

Example request bodies:

 POST /products
```json
{"article_code":"GLASS-001","name":"Glass"}
```
// POST /inbounds
```json
{"product_id":1,"quantity":10}
```

// POST /orders
```json
{"items":[{"product_id":1,"quantity":4}]}
```

Example end-to-end flow:

```powershell
curl -X POST http://localhost:8081/products -H "Content-Type: application/json" -d '{"article_code":"GLASS-001","name":"Glass"}'
curl -X POST http://localhost:8081/inbounds -H "Content-Type: application/json" -d '{"product_id":1,"quantity":10}'
curl http://localhost:8081/stock
curl -X POST http://localhost:8081/orders -H "Content-Type: application/json" -d '{"items":[{"product_id":1,"quantity":4}]}'
curl -X POST http://localhost:8081/orders/1/ship
curl http://localhost:8081/stock
```

## Business rules

- Product `articleCode` values are unique in PostgreSQL. Product names and article codes are trimmed and must be non-empty.
- Inbound quantities must be greater than zero, and the product must exist.
- Order item quantities must be greater than zero and an order must contain at least one item.
- Duplicate product IDs in one order are merged before stock is checked. For example, quantities `6` and `6` require 12 units.
- Creating an order does not reserve or subtract stock. Stock may change before shipment.
- Shipping locks and checks the order and stock in a transaction. If any item is short, no item is deducted and the order remains unchanged.
- Orders have two statuses: `CREATED` and `SHIPPED`. A shipped order cannot be shipped again.
- Stock is constrained to be non-negative by the database schema.

## Development

Run all tests:

```powershell
go test ./...
```

Run tests with coverage:

```powershell
go test ./... -cover
go test -coverprofile coverage.out ./...
go tool cover -func coverage.out
go tool cover -html coverage.out
```

Production code deliberately defines one interface: `order.Store` in `internal/order/service.go`. The order service uses it because order creation and shipping coordinate several storage operations and need focused business-rule tests. Product and stock code use concrete repositories directly; their tests use a mocked SQL connection. This keeps the project useful for learning when an interface helps instead of adding one for every struct.

Repository code owns SQL and transactions. Services own warehouse rules. Handlers own JSON binding and HTTP responses. The test suite covers more than 50% of production statements.

## Project layout

```text
cmd/api/main.go              Application entry point and route registration
internal/config/              Environment-variable configuration
internal/db/                  PostgreSQL connection setup
internal/product/             Product model, repository, and handler
internal/stock/               Inbound, stock report, service, repository, tests
internal/order/               Order, shipment, service, repository, tests
migrations/001_init.sql       PostgreSQL schema and goose directives
docker-compose.yml            Local PostgreSQL service definition
go.mod, go.sum                Go module metadata and dependencies
```

## Podman

```powershell
podman machine start podman-machine-default 
podman start mini_wms_postgres
```

## Migrations

Apply all pending migrations:

```powershell
go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable" up
```

Roll back the latest migration:

```powershell
go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations postgres "postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable" down
```

Create a migration:

```powershell
go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations create <name> sql
```

## Troubleshooting

### Windows: Port binding error on port 8081

**Symptom**: Application fails to start with error:
```
listen tcp :8081: bind: An attempt was made to access a socket in a way forbidden by its access permissions.
```
or
```
listen tcp :8081: bind: Only one usage of each socket address (protocol/network address/port) is normally permitted.
```

**Root cause**: 
1. **First error**: Windows reserves dynamic port ranges for Hyper-V, WSL, and other virtualization services. Port 8081 may fall within a reserved range (typically 7984-8083), preventing applications from binding to it.
2. **Second error**: Another process is already using port 8081.

**Fix for reserved port range** (requires Administrator privileges):

```powershell
# Stop WinNAT service
net stop winnat

# Exclude port 8081 from dynamic range
netsh int ipv4 add excludedportrange protocol=tcp startport=8081 numberofports=1

# Restart WinNAT service
net start winnat

# Verify exclusion (port 8081 should appear with asterisk)
netsh interface ipv4 show excludedportrange protocol=tcp
```

Port 8081 should now appear in the exclusion list with an asterisk (*). This exclusion persists across reboots.

**Fix for port already in use**:

```powershell
# Find and stop process using port 8081
Get-NetTCPConnection -LocalPort 8081 -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess | ForEach-Object { Stop-Process -Id $_ -Force }

# Verify port is free (should return nothing)
Get-NetTCPConnection -LocalPort 8081 -ErrorAction SilentlyContinue
```

## Out of scope

The current implementation does not provide a frontend, user accounts, auth, categories, explicit reservations, partial shipments, cancellation, pagination, filtering, multiple warehouses, a seed-data command, or generated OpenAPI/JSON Schema documentation.
