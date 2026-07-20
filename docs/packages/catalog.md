# Catalog Module

## Responsibility

`internal/catalog` will own product catalog behavior: defining products and making product information available to other use cases.

## Owned business concepts

- Product identity.
- Product descriptions needed to recognize and list products.

## Known use cases

- Create a product.
- Retrieve a product.
- List products.

## Belongs here

- Product-related business rules once they are known.
- Product lookup behavior required by catalog use cases.

## Does not belong here

- Stock quantities or availability.
- Receiving stock.
- Order lifecycle decisions.
- Shared technical infrastructure.

## Currently allowed dependencies

- Go standard library.
- Future internal packages only when a concrete use case justifies them.

## Currently forbidden dependencies

- HTTP frameworks such as Gin-Gonic.
- PostgreSQL drivers such as PGX.
- Another module's persistence implementation.
- Generic dumping-ground packages such as `common`, `shared`, `models`, or `utils`.

## Unresolved decisions

- Required product attributes.
- Whether product data can be changed or removed after use in an order.
- How product identity is represented externally and internally.

