# Dependency Rules

These rules describe the initial direction. They should stay small until real use cases require more detail.

## Established rules

- Application code belongs under `internal/`.
- Business modules are organized by capability: `catalog`, `inventory`, and `ordering`.
- `platform` is for shared technical infrastructure only and must not contain business logic.
- Business modules must not depend directly on HTTP frameworks, PostgreSQL drivers, or another module's persistence implementation.
- Avoid generic dumping-ground packages such as `common`, `shared`, `models`, and `utils`.
- Package boundaries should be introduced because of present responsibilities, not hypothetical future requirements.

## Current module distinctions

- `catalog` defines what a product is.
- `inventory` manages stock and availability.
- `ordering` owns the order lifecycle.
- Creating a catalog product must not automatically imply receiving stock.
- `ordering` should not own inventory quantities.
- `inventory` should not own product descriptions or the order lifecycle.

## Not decided yet

- Whether modules communicate through direct calls, interfaces, events, or another mechanism.
- Whether repository interfaces are needed, where they should live, or what they should contain.
- How database transactions will be coordinated across use cases.
- Which package names below each module are justified by actual code.

