# FUTO Home app template

* Go API, SQLite schema and tests. Read [server/AGENTS.md](server/AGENTS.md).
* SvelteKit routes, browser behavior and API consumers. Read [web/AGENTS.md](web/AGENTS.md).
* Setup and regeneration commands live in [README.md](README.md).

Keep each change focused on one outcome. For a behavior change, first write a test that shows the user visible failure at the cheapest tier that can catch it. Run the test before and after the change. Keep DB schema and OpenAPI artifacts in the same change as their sources.

Run `mise run check` from the root before handing off. For a runtime change, also exercise the affected flow against the running app. Use disposable data for schema and persistence checks. Report what passed and any cases you could not exercise.
