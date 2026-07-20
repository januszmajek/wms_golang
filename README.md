# WMS GOLANG

WMS GOLANG is a learning project for building a Go-based Warehouse Management System as a backend-only application.

The repository is currently an initial skeleton. It documents the product direction and architecture, but it does not yet implement WMS business behavior, HTTP handlers, persistence, migrations, or tests.

## Current architectural direction

- Modular monolith.
- Application code under `internal/`.
- Top-level business modules organized by capability: `catalog`, `inventory`, and `ordering`.
- `internal/platform` is reserved for genuinely shared technical infrastructure.
- Development should proceed through small backend vertical slices.
- PostgreSQL persistence is planned with PGX as the driver and goose for database migrations. Gin-Gonic is selected for HTTP requests and middleware. OpenAPI and JSON Schema are selected for validating requests/messages against their specifications. The logging library, code-generation approach, repository pattern, transaction strategy, and detailed package structure have not been selected.

## Repository map

```text
cmd/api/                 Minimal application entry point.
internal/catalog/        Future catalog module: products and product lookup.
internal/inventory/      Future inventory module: stock and availability.
internal/ordering/       Future ordering module: order lifecycle.
internal/platform/       Future shared technical infrastructure only.
migrations/              Future goose database migrations.
docs/product/            Product vision, use cases, and open questions.
docs/architecture/       Architecture overview and dependency rules.
docs/packages/           Business module responsibility notes.
docs/decisions/          Architecture Decision Records.
docs/learning/           Learning roadmap and progress notes.
```

## Commands

```powershell
go run ./cmd/api
go test ./...
```

`docker-compose.yml` currently provides a local PostgreSQL container for future persistence work. The Go application does not connect to it yet.

## Documentation starting points

- Product requirements: `docs/product/vision.md` and `docs/product/use-cases.md`
- Architecture overview: `docs/architecture/overview.md`
- Dependency rules: `docs/architecture/dependency-rules.md`
- Package map: `docs/architecture/package-map.md`
- Learning roadmap: `docs/learning/roadmap.md`

## Tooling status

Selected:

- Gin-Gonic for HTTP requests and middleware.
- PGX as the PostgreSQL driver.
- goose for database migrations.
- OpenAPI and JSON Schema for request/message specification and validation.

Not selected yet:

- Code-generation approach.
- Logger.
- Repository pattern.
- Transaction strategy.
- Detailed OpenAPI/JSON Schema validation workflow and library choices beyond the specifications themselves.
- Detailed package structure below the module level.
