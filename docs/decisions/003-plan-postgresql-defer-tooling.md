# 003: Plan PostgreSQL with PGX, goose Migrations, and JSON Schema

## Status

Accepted

## Context

Warehouse data will eventually need durable persistence. PostgreSQL is planned for that role. The project will use PGX as the PostgreSQL driver, goose for database migrations, and JSON Schema for schema definitions where JSON payload/schema documentation is needed.

The data-access pattern, transaction approach, JSON Schema file layout, supporting validation libraries, and any code-generation use have not been selected yet.

## Decision

Plan for PostgreSQL persistence. Use PGX as the PostgreSQL driver. Use goose as the migration tool. Use JSON Schema when the project needs explicit JSON payload/schema definitions.

Defer the data-access design, transaction strategy, and detailed JSON Schema workflow until a focused slice requires them.

## Consequences

- The repository can include `migrations/` as the location for goose migrations.
- The Go module can include PGX as the selected PostgreSQL driver dependency.
- No database access code, goose installation command, or local migration invocation convention is added now.
- Documentation may refer to goose as the selected migration tool and JSON Schema as the selected schema format for JSON payload/schema definitions.
- Documentation must avoid claiming that lib/pq, code generation, or any other database-adjacent tool has been selected.
- Early slices may use in-memory persistence if database persistence is not the learning objective.

## Unresolved implications

- Transaction management strategy.
- Data-access pattern.
- JSON Schema file layout, supporting validation libraries, and any code-generation use.

