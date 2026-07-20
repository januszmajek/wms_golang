# 002: Organize by Business Capability

## Status

Accepted

## Context

The known WMS areas are product catalog behavior, inventory behavior, and order lifecycle behavior. The project is intended to teach backend design through vertical slices, not through speculative technical layers.

## Decision

Organize the initial application code under business-capability modules:

- `internal/catalog`
- `internal/inventory`
- `internal/ordering`

Reserve `internal/platform` for shared technical infrastructure when a concrete need appears.

## Consequences

- Product, stock, and order lifecycle responsibilities are separated from the beginning.
- Creating a product does not automatically imply receiving stock.
- Technical subpackages should not be created until actual use cases require them.
- Generic packages such as `common`, `shared`, `models`, and `utils` are avoided.

## Unresolved implications

- Detailed package names below each module are not decided.
- The way modules collaborate is not decided.
- Aggregate boundaries and repository contracts are not decided.

