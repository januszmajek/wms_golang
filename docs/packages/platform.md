# Platform Module

## Responsibility

`internal/platform` is reserved for shared technical infrastructure when a real use case needs it.

## Owned concepts

No business concepts are owned by `platform`.

## Known use cases

None yet. Add code here only after a concrete technical responsibility is identified.

## Belongs here

Examples may eventually include technical adapters or process-wide infrastructure for selected tools such as Gin-Gonic, PGX, goose, OpenAPI, or JSON Schema, but only after a concrete integration need exists.

## Does not belong here

- Product, inventory, or ordering business logic.
- Reusable business rules.
- Generic helpers created before a concrete need exists.
- A substitute for `common`, `shared`, `models`, or `utils`.

## Currently allowed dependencies

- Go standard library.
- Selected infrastructure dependencies only when concrete code needs them.

## Currently forbidden dependencies

- Business logic from `catalog`, `inventory`, or `ordering`.
- Speculative wrappers around frameworks, databases, logging, or configuration.

## Unresolved decisions

- Which technical infrastructure, if any, should be centralized.
- How selected tools will be exposed to business modules without leaking framework or driver details.

