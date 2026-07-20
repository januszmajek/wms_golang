# 001: Choose a Modular Monolith

## Status

Accepted

## Context

WMS GOLANG is a learning project with a small set of related warehouse workflows. The project needs clear boundaries without the operational complexity of distributed services.

## Decision

Build WMS GOLANG as a modular monolith.

## Consequences

- The system can be developed and run as one Go application.
- Module boundaries can still be documented and enforced in code review.
- Cross-module behavior can evolve without network boundaries at the start.
- Poor package boundaries can still lead to coupling, so dependency rules must stay visible.

## Unresolved implications

- The exact mechanisms for module interaction are not selected.
- Transaction and persistence boundaries are not selected.

