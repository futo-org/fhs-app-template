# Server — Ubiquitous Language

Glossary for the `apps/templates/server` template. Terms here are the canonical
vocabulary; if code or conversation drifts from these, reconcile it.

## Entity

A Bun-tagged Go struct that maps 1:1 to a database table (e.g. `entity.User`
with `bun:"table:users,alias:u"`). Entities are the **source of truth for the
database schema**: Atlas generates migrations by diffing the entities against
the migration directory.

Entities are not confined to the repository — after a query they flow upward
through the service and are mapped to a **DTO** at the API boundary
(Immich-style). What the repository *contains* is the Bun query-building, not
the entity types themselves.

## DTO

A struct defining the shape of data crossing the **API boundary** — request
bodies and response payloads. Controllers accept and return DTOs. A DTO is
produced by mapping from one or more Entities (and may omit, rename, or combine
their fields).

## Domain Model — deliberately absent

There is **no** intermediate "pure domain" type between Entity and DTO. The
project maps `Entity → DTO` directly. A pure domain model is reintroduced only
if real domain behaviour (invariants, state machines) appears that does not
belong on either the Entity or the DTO.

## Cross-cutting errors

Sentinel errors that callers branch on (e.g. "not found", "invalid
credentials") are part of the shared language and live in a small dedicated
errors package, independent of any Entity or DTO.
