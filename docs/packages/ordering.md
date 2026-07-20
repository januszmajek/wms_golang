# Ordering Module

## Responsibility

`internal/ordering` will own the order lifecycle.

## Owned business concepts

- Order identity.
- Order items.
- Order lifecycle state once states are defined.
- Cancellation, shipping, and archiving behavior.

## Known use cases

- Create an order.
- Retrieve an order.
- List orders.
- Cancel an order manually.
- Ship an order.
- Archive an order.

## Belongs here

- Order lifecycle rules.
- Order item rules once defined.
- Coordination with inventory availability through an approved boundary when that boundary is designed.

## Does not belong here

- Inventory quantities.
- Product descriptions.
- Direct access to another module's persistence implementation.
- Generic workflow utilities.

## Currently allowed dependencies

- Go standard library.
- Future internal packages only when a concrete use case justifies them.

## Currently forbidden dependencies

- HTTP frameworks such as Gin-Gonic.
- PostgreSQL drivers such as PGX.
- Another module's persistence implementation.
- Generic dumping-ground packages such as `common`, `shared`, `models`, or `utils`.

## Unresolved decisions

- Whether order creation reserves stock.
- Whether order items are editable after order creation.
- Whether an empty order can be created.
- Whether a shipped order can be cancelled.
- Whether archiving is a lifecycle status or metadata such as `archived_at`.

