# Server tests

Run `mise run test` from the repository root. It uses the Ginkgo version from
`server/go.mod` and enables the race detector.

Use the cheapest test that exercises the behavior:

* Use a service with hand written repository fakes for business rules.
* Use `httptest` with the Fuego router for request validation and response mapping.
* Use `test.Configure()` for query and persistence checks. Each spec gets a fresh
  in memory SQLite database with the real migrations and repositories. Register
  its teardown with `DeferCleanup` and create rows through the app's services.

See the [album service spec](specs/services/album_service_test.go) for a database
test.
