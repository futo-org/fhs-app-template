# FUTO Home app template

Tasks run through mise from the repository root. [README.md](README.md) covers setup, checks, regeneration and production builds.

## Agent rules

1. This file holds the project layout, commands and rules.
2. Read [docs/agents/core.md](docs/agents/core.md) before writing code for the shared workflow, code shape, tests and commits.
3. Read the `AGENTS.md` for each component you touch.

The component file wins within its scope, then this file, then the shared rules. Claude Code loads the same rules through `CLAUDE.md` imports.

## Layout

| Area | Stack | Agent rules |
| --- | --- | --- |
| `server/` | Go API with Fx, Fuego, Bun and SQLite | [Go rules](server/AGENTS.md) |
| `web/` | Static SvelteKit app | [SvelteKit rules](web/AGENTS.md) |
| `scripts/` | Release scripts | [Shell rules](scripts/AGENTS.md) |

The shell rules also apply to scripts in mise tasks and GitHub workflows. Compose files are in `docker/`. `mise run dev` starts the server and web app. `mise run dev-down` stops them.

## Checklist

Run `mise run check` for the same checks as CI. The component guides list the focused commands. Exercise runtime changes against a disposable running app and report the states you checked.

| Changed | Regenerate |
| --- | --- |
| Server route or DTO | `mise run open-api`, then commit `server/doc/openapi.json` |
| Bun entity | `mise run db:diff <name>`, then commit the SQL and `atlas.sum` |
| Handwritten migration | `mise run db:hash`, then commit the SQL and `atlas.sum` |

## Where the knowledge lives

Before starting, look in `docs/` for the design, implementation plan and test plan for the area. Read the relevant files, the component README and its `CONTEXT.md` when present. The server vocabulary lives in [server/CONTEXT.md](server/CONTEXT.md). Bring discrepancies between code and design to the developer.

When adding project documentation, put designs in `docs/design/`, milestone plans in `docs/implementation_plans/`, and manual test plans in `docs/testing/`. Keep the plan's task IDs and links stable. A test plan gives setup pitfalls, the steps and their expected results, and any hardware requirements.

## Shared review expectations

Read the [review requirements](docs/review/requirements.md) before reviewing changes or answering a review thread.

## Code comments

Do not add code comments. Make the code self explanatory through naming and structure.

## Go services

* Keep `server/internal/service` for `*_service.go`, `module.go` and tests. Put methods in their service file.
* Put DTOs and other logic that belongs outside the service in the appropriate package. Check the import graph before moving code.
* Split large services by domain. Give each its own type, interface and file. Aim for 100 to 300 lines per file and split files approaching 500 lines.

## UI components

Reuse the app's components and design tokens. If a UI kit is installed, check the pinned version before using a component. When a new component is needed, name it in the PR and explain the gap it fills.

## End to end tests

* For apps with an end to end suite, cover new user flows and run the affected specs against the dev stack.
* Find elements by role, label and text. Assert what the user sees. Keep CSS selectors in a shared page helper when needed.
* Seed through the API with disposable data for the run. Keep test media small.

## Commits

Use the shared commit conventions. When a scope helps, use the top level directory or package.
