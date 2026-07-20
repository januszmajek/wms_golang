# Inventory Module

## Responsibility

`internal/inventory` will manage stock and availability for products that have been registered in inventory.

## Owned business concepts

- Inventory registration.
- Stock quantity.
- Availability rules once defined.

## Known use cases

- Register a product in inventory.
- Receive stock to increase quantity.
- Retrieve stock information.
- Future reservation and release of stock.

## Belongs here

- Rules for receiving stock.
- Rules for reporting stock information.
- Availability calculations once product requirements define them.

## Does not belong here

- Product descriptions owned by `catalog`.
- Order lifecycle behavior owned by `ordering`.
- Automatic product creation.
- Generic shared business helpers.

## Currently allowed dependencies

- Go standard library.
- Future internal packages only when a concrete use case justifies them.

## Currently forbidden dependencies

- HTTP frameworks such as Gin-Gonic.
- PostgreSQL drivers such as PGX.
- Another module's persistence implementation.
- Generic dumping-ground packages such as `common`, `shared`, `models`, or `utils`.

## Unresolved decisions

- Whether there is one warehouse or multiple warehouses.
- Whether stock tracks only total quantity or separate available and reserved quantities.
- Whether order creation reserves stock.
- When available stock is reduced.

