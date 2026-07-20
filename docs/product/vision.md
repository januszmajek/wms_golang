# Product Vision

Mini WMS is a small backend Warehouse Management System for learning Go, backend design, and modular monolith architecture.

## Current goal

Provide a JSON API for core warehouse workflows around products, inventory, and orders. Gin-Gonic is selected for HTTP requests and middleware. OpenAPI and JSON Schema are selected to validate requests/messages against their specifications. PostgreSQL persistence is planned with PGX as the driver and goose for migrations.

## Known capabilities

- Create a product.
- Retrieve a product.
- List products.
- Register a product in inventory.
- Receive stock to increase quantity.
- Retrieve stock information.
- Create an order.
- Retrieve an order.
- List orders.
- Cancel an order manually.
- Ship an order.
- Archive an order.

## Current constraints

- Backend-only; no frontend is planned at this stage.
- PostgreSQL persistence is planned with PGX as the driver and goose for database migrations.
- The first implementation slices should stay small enough to support learning and review.

## Open product questions

- Does the system support one warehouse or multiple warehouses? One
- Does an order reserve stock before shipping? Yes
- At which point is available stock reduced? When order is created
- Is archiving an order a lifecycle status or separate metadata such as `archived_at`? status
- Can a product be changed or removed after it has been used in an order? No
- Are order items editable after order creation? No
- Can an empty order be created? No
- Can a shipped order be canceled? No
- Will the inventory initially track only total quantity, or separate available and reserved quantities? One quantity
- Which OpenAPI/JSON Schema validation workflow and supporting Go libraries should be used?
- Will OpenAPI or JSON Schema be used for code generation later?

