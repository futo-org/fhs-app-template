# Go server

* Controllers register named Fuego handlers. They parse and validate input, call a service, and map errors to HTTP at the transport edge.
* Services own business rules and map Bun entities to DTOs. Inject repository interfaces through constructors and keep SQL in repositories. Keep Fx constructors free of new I/O. Manage startup, shutdown and background work through the lifecycle.
* Use request contexts through every layer. Return `apperr` sentinels with `%w` context. Match with `errors.Is` or `errors.As` and log failures at the edge.
* Keep request and response shapes in `internal/dto`. When a route or DTO changes, run `mise run open-api` and commit the spec. The export registers the real routes and exits without running migrations or listening.
* For entity changes, read the [migration guide](internal/database/migrations/README.md). Run `mise run db:diff <name>` to generate the migration and commit it with the checksum. Exercise the migration from existing data when changing the shape of stored data.
* Use Ginkgo and Gomega. Use injected repository fakes for service rules, `httptest` for HTTP validation and status codes, and the [SQLite harness](test/README.md) for queries and cross table effects. Label database specs `medium` and isolated specs `unit`. Assert the specific sentinel on failure.
