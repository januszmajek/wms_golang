# 005: Use OpenAPI and JSON Schema for Message Validation

## Status

Accepted

## Context

The project will expose JSON messages. The project needs a way to describe expected request/message shapes and validate whether messages align with their specifications.

## Decision

Use OpenAPI and JSON Schema to specify and validate HTTP requests/messages.

## Consequences

- API payloads should eventually have specifications that can be checked during request/message validation.
- Validation against message specifications is separate from business-rule validation.
- Documentation may refer to OpenAPI and JSON Schema as selected specification formats.
- The project should not invent request schemas before the related vertical slice defines the message.

## Unresolved implications

- Exact OpenAPI file layout and naming.
- Exact JSON Schema file layout and naming.
- Go libraries or tools used to perform validation.
- Whether OpenAPI or JSON Schema will later be used for code generation.
- How validation errors are exposed in API responses.

