# Architecture Overview

WMS GOLANG is planned as a modular monolith.

## Established decisions

- Application code belongs under `internal/`.
- Business capabilities are the primary top-level division.
- The initial business modules are `catalog`, `inventory`, and `ordering`.
- `internal/platform` is reserved for genuinely shared technical infrastructure and must not contain business logic.
- The application will be developed through small end-to-end backend vertical slices.
- A vertical slice can include HTTP transport, application behavior, domain rules, persistence, and tests. It does not require a frontend.
- PostgreSQL persistence is planned with PGX as the driver and goose for database migrations.
- Gin-Gonic is selected for HTTP requests and middleware.
- OpenAPI and JSON Schema are selected for validating requests/messages against their specifications.

## Current repository shape

```text
cmd/api/            Application entry point.
internal/catalog/   Product catalog capability.
internal/inventory/ Stock and availability capability.
internal/ordering/  Order lifecycle capability.
internal/platform/  Shared technical infrastructure when justified.
docs/               Product, architecture, package, decision, and learning notes.
migrations/         Future goose database migrations.
```

## Deliberately deferred

- Transaction strategy.
- Code-generation tooling.
- Detailed OpenAPI/JSON Schema validation workflow and supporting libraries.
- Detailed aggregate boundaries.
- Repository interfaces and persistence package layout.
- Cross-module integration mechanisms.

