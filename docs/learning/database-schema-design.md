# Database Schema Design

**Date:** 2026-07-21 to 2026-07-22  
**Learning Focus:** Designing a complete warehouse management database schema through migrations  
**Status:** Implemented and tested

---

## Schema Overview

The warehouse database consists of seven tables, organized across three business domains:

### Products (Catalog Domain)
```
products
├── id (SERIAL PRIMARY KEY)
├── product_code (VARCHAR(50) UNIQUE NOT NULL)
└── name (VARCHAR(75) NOT NULL)
```

**Rationale:**
- `id`: Surrogate key for database relationships (easier to maintain than natural keys)
- `product_code`: Business identifier customers use (UNIQUE to prevent duplicates)
- `name`: Human-readable product name

### Inventory (Inventory Domain)
```
stock
├── product_id (BIGINT PRIMARY KEY → products)
├── quantity (INTEGER NOT NULL CHECK >= 0)
└── updated_at (TIMESTAMP DEFAULT NOW())

inbound_operations
├── id (BIGSERIAL PRIMARY KEY)
├── product_id (BIGINT NOT NULL → products)
├── quantity (INTEGER NOT NULL CHECK > 0)
└── created_at (TIMESTAMP DEFAULT NOW())
```

**Rationale:**
- `stock`: Single row per product tracking current quantity
- `inbound_operations`: Audit trail of received inventory
- CHECK constraints ensure valid quantities (stock >= 0, inbound > 0)

### Orders (Ordering Domain)
```
orders
├── id (BIGSERIAL PRIMARY KEY)
├── status (TEXT NOT NULL)
├── created_at (TIMESTAMP DEFAULT NOW())
└── shipped_at (TIMESTAMP NULL)

order_items
├── id (BIGSERIAL PRIMARY KEY)
├── order_id (BIGINT NOT NULL → orders)
├── product_id (BIGINT NOT NULL → products)
└── quantity (INTEGER NOT NULL CHECK > 0)

outbound_operations
├── id (BIGSERIAL PRIMARY KEY)
├── order_id (BIGINT NOT NULL UNIQUE → orders)
└── created_at (TIMESTAMP DEFAULT NOW())
```

**Rationale:**
- `orders`: Order metadata and lifecycle (status, creation, shipment timestamp)
- `order_items`: Line items in an order (references both order and product)
- `outbound_operations`: Shipment event tracking
- `UNIQUE` on `outbound_operations.order_id`: Enforces business rule that orders ship completely (no partial shipments)

---

## Key Design Decisions

### 1. Surrogate Keys (id columns)

All tables use auto-generated numeric IDs as primary keys, not business values.

**Why:**
- Stable: If product_code changes, foreign keys don't break
- Efficient: Integers are smaller than strings for joins
- Future-proof: Easy to support distributed systems later

### 2. Auditability

Both `inbound_operations` and `outbound_operations` track **when** things happened:
- `inbound_operations.created_at`: When stock was received
- `outbound_operations.created_at`: When orders were shipped

This creates an audit trail of inventory movements.

### 3. Quantity Constraints

- `stock.quantity >= 0`: Can't go negative
- `inbound_operations.quantity > 0`: Only positive receipts make sense
- `order_items.quantity > 0`: Only positive order quantities

These CHECK constraints prevent invalid data at the database level.

### 4. Complete Order Shipment

`outbound_operations.order_id` is UNIQUE, enforcing: **one outbound operation per order**.

Business rule: Orders are always shipped completely, never partially.

---

## Relationships

```
products
  ├→ stock (1:1)
  ├→ inbound_operations (1:many)
  ├→ order_items (1:many)
  └→ orders (indirect, via order_items)

orders
  ├→ order_items (1:many)
  └→ outbound_operations (1:1)
```

**Full workflow:**
1. Create product in `products`
2. Receive inventory → `inbound_operations` + update `stock.quantity`
3. Create order → `orders` + `order_items` (references products)
4. Ship order → `outbound_operations` (references orders)

---

## What This Teaches

### Concepts Practiced
- Table relationships (foreign keys, 1:1, 1:many)
- Constraints (PRIMARY KEY, UNIQUE, CHECK, REFERENCES)
- Audit trail patterns (created_at, updated_at timestamps)
- Business rule enforcement at database level
- Type consistency (BIGINT vs SERIAL)

### Design Principles
- Separate concerns: Products, inventory, orders are distinct
- Immutability: Inbound/outbound operations are never deleted (audit trail)
- Validation: Constraints prevent invalid states
- Traceability: Timestamps and IDs enable debugging

---

## What Comes Next

1. **Insert product into database** (CREATE, READ operations)
2. **HTTP endpoint to list/create products** (transport layer)
3. **JSON Schema for request/response** (specification)
4. **Quicktype to generate Go structs** (learning the automation target)
5. **Complete order-to-shipment workflow** (implement use cases)

---

## Running the Migrations

```powershell
# Check status
goose -dir migrations postgres "$env:DATABASE_URL" status

# Apply all pending migrations
goose -dir migrations postgres "$env:DATABASE_URL" up

# Verify tables exist
docker exec -it postgres-db psql -U admin -d wms_golang_db -c "\dt"
```

**Files:**
- `migrations/20260721091402_add_products_table.sql` — Products table
- `migrations/20260721201805_other_tables.sql` — All other tables (stock, operations, orders)

