# Copilot Instructions — WMS GOLANG

## Project overview

WMS GOLANG is a Go learning project for a backend-only Warehouse Management System. The repository is currently an initial skeleton: documentation and module directories exist, but WMS behavior, HTTP handlers, persistence, migrations, and tests are not implemented yet.

## Commands

```powershell
# Run the minimal entry point
go run ./cmd/api

# Run all current Go checks/tests
go test ./...

# Start the local PostgreSQL container for future persistence work
docker compose up -d
```

Do not add dependency installation commands unless a selected tool is needed by code or documented project workflow. goose is selected for migrations, but do not add migration commands until migration files and the local invocation convention are documented.

## Architecture

- The system is a modular monolith.
- Application code belongs under `internal/`.
- Business capabilities are the primary top-level division.
- Initial business modules are `catalog`, `inventory`, and `ordering`.
- `internal/platform` is reserved for genuinely shared technical infrastructure and must not contain business logic.
- Development should proceed through small backend vertical slices.
- A vertical slice can include HTTP transport, application behavior, domain rules, persistence, and tests. It does not require a frontend.

## Current structure

```text
cmd/api/main.go
internal/catalog/
internal/inventory/
internal/ordering/
internal/platform/
migrations/
docs/product/
docs/architecture/
docs/packages/
docs/decisions/
docs/learning/
go.mod
```

## Dependency and package rules

- Do not create speculative subpackages such as `domain`, `application`, `http`, `postgres`, `repositories`, `services`, `models`, `shared`, `common`, or `utils`.
- Introduce package boundaries because of present responsibilities, not hypothetical future requirements.
- Business modules must not depend directly on HTTP frameworks, PostgreSQL drivers, or another module's persistence implementation.
- `catalog` defines what a product is.
- `inventory` manages stock and availability.
- `ordering` owns the order lifecycle.
- Creating a catalog product must not automatically imply receiving stock.
- `ordering` should not own inventory quantities.
- `inventory` should not own product descriptions or the order lifecycle.

## Tooling status

PostgreSQL and migrations are planned. PGX is selected as the PostgreSQL driver, goose is selected as the migration tool, Gin-Gonic is selected for HTTP requests and middleware, and OpenAPI plus JSON Schema are selected for request/message specification and validation.

The logging library, repository abstraction, transaction approach, code-generation approach, and detailed OpenAPI/JSON Schema validation workflow have not been selected.

## Tutor and learning context

- Agent tutor rules live in `.github/skills/golang-tutor/SKILL.md`.
- Learning project notes live under `docs/learning/`.
- Use `README.md`, `docs/`, current code, and learning progress as sources of truth before mentoring or changing code.

## Communication with the user

Writing rules, from Orwell, 1946. These govern prose: docs, PR text, messages. Never touch code or technical terms;
swap in everyday words only where precision survives.

1. Never use a metaphor, simile or other figure of speech which you are used to seeing in print.
2. Never use a long word where a short one will do.
3. If it is possible to cut a word out, always cut it out.
4. Never use the passive where you can use the active.
5. Never use a foreign phrase, a scientific word or a jargon word if you can think of an everyday English equivalent.
6. Break any of these rules sooner than say anything outright barbarous.
   
Review every prose output against these rules before delivering.
