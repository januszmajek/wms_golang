# 004: Use Gin-Gonic for HTTP Requests and Middleware

## Status

Accepted

## Context

The project will expose an HTTP JSON API. The HTTP framework and middleware approach had previously been deferred.

## Decision

Use Gin-Gonic for HTTP requests and middleware.

## Consequences

- Future HTTP entry points and middleware can be built with Gin-Gonic.
- Business modules must still not depend directly on Gin-Gonic.
- HTTP-specific code should stay at the transport/infrastructure edge rather than inside domain behavior.
- The project can add Gin-Gonic dependency entries before the first HTTP slice, but it should not add handlers or middleware until a vertical slice requires them.

## Unresolved implications

- Exact route structure.
- Middleware list and ordering.
- Error response format.
- How request validation failures will be represented in API responses.

