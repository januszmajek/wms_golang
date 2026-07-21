# Database Migrations and PostgreSQL Setup

## Date
2026-07-21

## Learning Objective
Set up database migrations using goose and connect the Go application to PostgreSQL.

## Concepts Practiced

### 1. Database Migrations with Goose
- Installed goose migration tool: `go install github.com/pressly/goose/v3/cmd/goose@latest`
- Documented goose commands in `migrations/README.md`
- Created first migration: `20260721091402_add_products_table.sql`
- Applied migration using: `goose -dir migrations postgres "$DATABASE_URL" up`
- Verified migration status and rollback capability

### 2. PostgreSQL Table Design
- Designed `products` table structure based on use case requirements
- Made design decisions about keys:
  - **Surrogate key (`id`)**: Auto-incrementing SERIAL for primary key and foreign key relationships
  - **Natural key (`product_code`)**: Business-meaningful identifier with UNIQUE constraint
- Applied NOT NULL constraints to prevent invalid data
- Chose appropriate data types:
  - `id`: SERIAL PRIMARY KEY (auto-incrementing)
  - `product_code`: VARCHAR(50) UNIQUE NOT NULL
  - `name`: VARCHAR(75) NOT NULL

### 3. SQL in Go with PGX
- Used `pgxpool.Pool.QueryRow()` to execute SQL queries
- Applied `.Scan()` to map database column results to Go variables
- Followed the pattern:
  ```go
  var result int
  err := db.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&result)
  if err != nil {
      // handle error
  }
  ```

### 4. Design Trade-offs
- **Natural vs. Surrogate Keys**: Understood when to use business identifiers vs. technical IDs
- **Stability**: Surrogate keys protect foreign key relationships when business identifiers change
- **Simplicity**: Integer primary keys are faster for joins and simpler for debugging

## Key Decisions

### Migration File Location
Stored in `migrations/` directory at project root, as planned in ADR 003.

### Table Naming Convention
Used lowercase with underscores (e.g., `products`), following PostgreSQL conventions.

### Command Documentation
Documented complete goose commands with all required parameters:
- Database driver: `postgres`
- Connection string: from `DATABASE_URL` environment variable
- Migration directory: `-dir migrations`

## Files Created/Modified

### Created
- `migrations/20260721091402_add_products_table.sql` — First migration creating products table

### Modified
- `migrations/README.md` — Documented goose installation and usage commands
- `cmd/api/main.go` — Added query to count products table rows

## Next Steps

### Immediate
- Create HTTP endpoint to query products
- Define Go struct to represent Product
- Return products as JSON

### Future (JSON Schema Workflow)
- Define data models in JSON Schema
- Use quicktype.io to generate Go structs
- Replace manual struct definitions with generated code
- This aligns with target workplace practices

## Questions Resolved

**Q: Should I use JSON Schema to create database tables?**
A: No. JSON Schema validates JSON payloads (HTTP requests/responses). SQL migrations create database tables. They solve different problems and both should coexist.

**Q: Do I need both `id` and `product_code`?**
A: Yes. Use `id` as the surrogate primary key for database relationships. Use `product_code` as the unique business identifier. This separation provides stability when business values change.

**Q: Do I have to write SQL in Go code?**
A: For now, yes. Writing SQL strings in Go is the simplest approach for learning. Later, you can explore query builders, ORMs, or code generation tools like sqlc.

## Learning Reflection

This session established the foundation for database persistence:
- Migrations manage schema changes over time
- SQL strings in Go code execute queries and map results
- Design decisions (keys, constraints, types) depend on use case requirements, not abstract principles

The next learning objective is HTTP request handling, where JSON Schema becomes relevant for generating request/response structs.

