# FUTO Home app template

This starter pairs a Go API and SQLite database with a SvelteKit web app. The server serves the static web build and applies embedded migrations in production.

## Develop

Install [mise](https://mise.jdx.dev/) and Docker with Compose. From the repository root:

```sh
mise trust
mise install
mise run dev
```

The web app runs at [localhost:3000](http://localhost:3000) and the API docs at [localhost:3001/swagger](http://localhost:3001/swagger). Run `mise run dev-down` to stop both.

## Check

Run `mise run check` from the repository root. CI runs the same task. It runs golangci-lint and the Go tests with the race detector, builds the Go binary and web app, checks Svelte types, and compares the OpenAPI spec and migrations to their Go sources. Running the Go tests with the race detector requires a C compiler. For Go lint alone, run `mise run server:lint`.

| Changed | Command | Commit with the change |
| --- | --- | --- |
| Go source formatting | `mise run server:format` | Formatted source |
| API route or DTO | `mise run open-api` | `server/doc/openapi.json` |
| Bun entity | `mise run db:diff <name>` | SQL migration and `atlas.sum` |
| Handwritten migration | `mise run db:hash` | SQL migration and `atlas.sum` |

See the [migration guide](server/internal/database/migrations/README.md) for schema changes and [agent rules](AGENTS.md) for project conventions.

## Build the production image

Run `mise run prod-build` to build the image and print its version. Run `mise run prod` to serve the web app and API on port 3001 and `mise run prod-down` to stop it.
