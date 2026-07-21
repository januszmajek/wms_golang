# Migrations


## Commands

Goose was installed on WSL using Go. It needs access to the migrations folder and the database connection string. 

In PowerShell, set the DATABASE_URL environment variable before running goose commands:
```powershell
$env:DATABASE_URL = "postgres://admin:123@<host>:5432/wms_golang_db?sslmode=disable"
```

Or load it from your `.env` file if using a tool that reads it.

### Installation

```powershell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Create a migration


```powershell
goose -dir migrations create <migration_name> sql
```

Goose automatically prefixes the filename with a timestamp (e.g., `20260721091402_add_products_table.sql`).

### Run migrations

To apply all pending migrations:

```powershell
goose -dir migrations postgres "$env:DATABASE_URL" up
```

### Rollback migrations

To rollback the most recent migration:

```powershell
goose -dir migrations postgres "$env:DATABASE_URL" down
```

### Check migration status


```powershell
goose -dir migrations postgres "$env:DATABASE_URL" status
```

