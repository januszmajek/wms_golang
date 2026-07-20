# Product Use Cases

This file records only the currently known WMS use cases. It does not define API routes, Go types, database tables, or business validation rules that have not yet been decided. OpenAPI and JSON Schema are selected for request/message specification validation.

## Catalog

1. Create a product.
2. Retrieve a product.
3. List products.

## Inventory

1. Register a product in inventory.
2. Increase stock quantity by receiving stock.
3. Retrieve stock information.

## Ordering

1. Create an order.
2. Retrieve an order.
3. List orders.
4. Cancel an order manually.
5. Ship an order.
6. Archive an order.

## Initial learning slice order

The provisional implementation order is:

1. Create a product.
2. Retrieve and list products.
3. Receive stock.
4. Retrieve stock information.
5. Create an order.
6. Cancel an order.
7. Ship an order.
8. Archive an order.

This order is guidance, not a fixed implementation plan. It can change when domain questions are answered.

