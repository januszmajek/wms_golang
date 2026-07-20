# Package Map

This map records intended responsibilities without inventing subpackages or public APIs.

## `cmd/api`

Minimal application entry point. It currently does not start a server. Gin-Gonic is selected for future HTTP requests and middleware.

## `internal/catalog`

Future business module for product creation and product lookup.

## `internal/inventory`

Future business module for inventory registration, receiving stock, and stock information.

## `internal/ordering`

Future business module for order creation, retrieval, listing, cancellation, shipping, and archiving.

## `internal/platform`

Reserved for genuinely shared technical infrastructure when a concrete need appears. It must not become a home for reusable business concepts.

## `migrations`

Reserved for future goose database migrations. PostgreSQL is planned with PGX as the driver, but detailed data-access design is not selected.

## `docs`

Repository-level product, architecture, decision, package, and learning documentation.

