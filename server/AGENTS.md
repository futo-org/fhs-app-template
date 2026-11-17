# Go rules · extends [core](../docs/agents/core.md)

Governs `server/`. See [server vocabulary](CONTEXT.md) for Bun entities, DTOs and shared errors.

## Skills to load

When the `samber/cc-skills-golang` pack is installed, load `samber/cc-skills-golang@golang-how-to` first and follow its routing. Use its style, naming, error handling, testing, troubleshooting, safety, security, design, data structures, documentation, modernization and CI guidance. This file sets project overrides.

## Vocabulary

| Core | Here |
| --- | --- |
| transport | Fuego handlers in `internal/controller` parse inputs and map service errors to HTTP |
| logic | services in `internal/service` apply business rules and map entities to DTOs |
| adapter | Bun repositories in `internal/repository` or a client package, one technology each |
| shape | request and response structs with JSON tags in `internal/dto` |
| `findOrFail` | read returning `(T, error)` wrapping `apperr.ErrNotFound` |
| `requireX` | returns an error when the condition fails |
| `checkX`, `isX`, `hasX` | predicate returning a bool or the allowed set |

## Layout

```text
main.go               Fx wiring
openapi.go            OpenAPI export from the route graph
internal/
  apperr/             shared sentinel errors
  controller/         HTTP input validation and responses
  database/entity/    Bun entities and desired schema
  database/migrations/ SQL migrations applied by Goose
  dto/                request and response types
  repository/         queries and I/O
  service/            business rules
test/
  setup.go            SQLite test setup
  specs/              Ginkgo suites
```

## Idioms

* Use guard clauses and early returns. Use an options struct for a function with more than four parameters and functional options for configurable constructors.
* Reuse the dependency's exported interface. When none exists, declare an unexported interface in the consuming package with the methods the caller uses.
* Inject dependencies through constructors. Register providers in Fx modules. Put listeners, connections and other I/O in lifecycle hooks and close them on stop.
* Use `map[T]struct{}` or an existing set type for membership. Use `errgroup` for independent I/O and `SetLimit` for fan out.
* Use `gofmt -r` for mechanical replacements across code where it fits.
* Use `Unknown` for enum zero values. Name constants by role in the package that owns them.
* Declare implementation and fake interface assertions, for example `var _ repository.AlbumRepository = (*FakeAlbums)(nil)`.
* Pass the request context through every layer. Use `context.WithoutCancel(ctx)` for work that must outlive the request and give background work an owner that can wait for it.

## Size

* Split files over 500 lines and functions over 80 lines.
* Register named handlers with Fuego. Keep `RegisterRoutes` for route declarations.

## Boundaries

* Inject a dependency to replace it in tests. Keep production behavior out of mutable package variables and wire registries with Fx providers.
* Use the caller's context in blocking work. Reserve fresh root contexts for entry points and lifecycle work.
* Pass configuration into logic. Read environment variables in `internal/config`.
* Return errors for failed requests. Keep `Must*` and panics in initialization.
* Define fakes beside tests. Keep direct console output in the CLI entry point and use `slog` for application logs.

## Lints

[.golangci.yml](.golangci.yml) enables standard golangci-lint checks plus `intrange`, `modernize`, `revive` and `gofmt`. Fix findings at the source. Record the reason for any necessary exclusion.

## Errors

```go
album, err := s.albums.GetByID(ctx, id)
if err != nil {
	return nil, fmt.Errorf("get album %d: %w", id, err)
}
```

* Return the appropriate `apperr` sentinel and wrap context with `%w`. Check with `errors.Is` or `errors.As`.
* Map shared sentinels in one transport handler. Keep detailed error chains in logs. Give a missing resource and a resource the caller cannot see the same response.
* Log once at the edge. Use `slog` with context when available, `Debug` for expected failures such as a client abort, and `Error` for unknown failures.
* Run goroutines through an owner such as an `errgroup` that waits and reports failures.

## Tests

Use Ginkgo and Gomega. Label isolated specs `unit` and real SQLite specs `medium`.

| Tier | Harness |
| --- | --- |
| Unit | service with injected repository fakes |
| Transport | `httptest` against the Fuego router with a fake service |
| Integration | [test.Configure](test/README.md) with fresh SQLite and the real migrations and repositories |

* Write fakes by hand. Embed the repository interface, implement the methods the test calls and record calls that matter.
* Assert sentinels with `MatchError(apperr.ErrX)`. Use `DescribeTable` and `Entry` for cases that differ by one input.
* Put shared factories and fakes under `test/`. Plain `testing` with `TestType_Behavior` names fits small binaries without a Ginkgo suite.

## Checklist

```sh
mise run server:format
mise run server:lint
mise run test
```

Run the SQLite integration tier after changing a query or access rule. See the [test guide](test/README.md) for setup.

## Derived artifacts

Use `mise run open-api` after changing a route or DTO. Export registers the real routes and exits without running migrations or listening. Commit [the spec](doc/openapi.json) with the source change.

Read the [migration guide](internal/database/migrations/README.md) before changing an entity or migration. Generate schema changes with `mise run db:diff <name>` and commit the SQL with `atlas.sum`. For transforms, exercise the migration from realistic existing data. CI checks the schema with `mise run db:check`.

## Overrides

* The Core `findOrFail` and `requireX` conventions return errors in Go. HTTP access policy belongs at the edge.
* Use package `slog` and the context variants when a context is available.
* Use the repository rules for code comments and Ginkgo tests when applying the Go skill pack. Keep dependency free shared functions in a small package named for their role.
