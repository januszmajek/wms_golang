1# Copilot Instructions — Mini WMS

## Project Overview

Mini WMS is a Go REST API for a small warehouse. Backend-only (no frontend), no auth. Test with `curl`.

## Commands

```bash
# Run server
go run ./cmd/api

# Run all tests
go test ./...

# Run a single package's tests
go test ./internal/order/...
go test ./internal/stock/...

# Run with coverage (must be >50%)
go test ./... -cover

# Generate HTML coverage report
go test ./... -coverprofile=cov.out && go tool cover -html=cov.out

# Start DB
docker compose up -d

# Apply migrations
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable" up

# Rollback migrations
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable" down

# Create a new migration
goose -dir migrations create <name> sql
```

## Architecture

```
HTTP handlers -> Services -> Repositories -> DB
```

Strict layering under `internal/`:

| Layer | Responsibility |
|---|---|
| **handlers** | Parse/validate JSON, call service, return HTTP status + JSON |
| **services** | Business logic only — no HTTP, no SQL. Unit-tested with mocks |
| **repositories** | SQL only — CRUD, stock updates, DB transactions |

Each domain (`product`, `stock`, `order`) has its own subdirectory with a `model.go`, `repository.go`, `service.go`, `handler.go`, and `service_test.go`.

## Project Structure

```
cmd/api/main.go
internal/
  config/config.go
  db/db.go
  product/{model,repository,handler}.go
  stock/{model,repository,service,handler,service_test}.go
  order/{model,repository,service,handler,service_test}.go
migrations/001_init.sql
docker-compose.yml
.env.example
go.mod
```

## Key Conventions

### Service testing
- Unit tests live in `service_test.go` within the domain package.
- Services are tested with mocked repositories (interface-based), not against a real DB.
- `stock` and `order` are the only domains with services (and tests); `product` has no service layer.

### Business rules to enforce in code
- **Inbound**: quantity must be > 0; product must exist.
- **Order creation**: validates stock *before* creating. Does NOT decrease stock. Status starts as `CREATED`. Duplicate `product_id` entries in a single request must be summed before checking stock.
- **Shipment** (`POST /orders/:id/ship`): runs in a DB transaction: check status → re-check stock → decrease stock → update order to `SHIPPED` → record outbound operation. Reject if already `SHIPPED` or stock is now insufficient.
- Stock quantity is always a non-negative integer (`CHECK (quantity >= 0)` in DB).
- No partial fulfillment; orders either fully succeed or fail.

### Order statuses
- `CREATED` — stock checked and sufficient, not yet decreased.
- `SHIPPED` — stock decreased, outbound operation recorded.

### Error responses
- HTTP 400 for business rule violations (insufficient stock, invalid input, already shipped).
- HTTP 404 for missing resources.
- HTTP 500 for unexpected DB errors.

### Environment config
Loaded from `.env` (see `.env.example`):
```
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable
```

### Database
- PostgreSQL via `database/sql` + `lib/pq` driver.
- Migrations managed with `goose` in `migrations/`.
- Shipment logic uses a DB transaction (begin/commit/rollback pattern in repository layer).
