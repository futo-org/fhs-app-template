# Database Migrations

This project uses a split workflow:

- **Atlas** *authors* migrations by diffing your Bun **entities** against the
  migration history. It runs on **dev/CI only** — never on a device.
- **goose** *applies* migrations at runtime, as a pure-Go library embedded in the
  binary (`migrations/module.go`, an `fx.Invoke`). No `goose`/`atlas` binary ships
  with the app.

```
internal/database/entity/*.go   ← source of truth for the schema
        │
        │  mise run db:diff <name>      (Atlas, dev/CI)
        ▼
internal/database/migrations/*.sql   ← generated goose-format files (committed)
        │
        │  goose.Up() at startup       (runtime, on-device, forward-only)
        ▼
      SQLite database
```

The **entities are the source of truth.** You do not hand-write schema SQL — you
edit a struct and let Atlas compute the migration.

## Prerequisites

The Atlas CLI is provided by mise (`atlas = "1.2.0"` in `../mise.toml`). Confirm it
resolves:

```bash
mise run --help >/dev/null && atlas version
```

`goose` is a Go library here (applied automatically at startup); the goose CLI is
only needed for the hand-written case below.

## Creating a migration from an entity change (the normal path)

1. Edit or add a struct in `internal/database/entity/` (e.g. add a column to
   `entity.Media`, or add a new `entity.Album`). Every entity must live in this one
   package — that's how Atlas's loader discovers them.

2. Generate the migration. The final argument is the migration name:

   ```bash
   mise run db:diff add_media_created_at
   ```

   Atlas spins up a throwaway in-memory SQLite dev database, loads the desired
   schema from your entities, diffs it against the files already in this directory,
   and writes a new timestamped `*.sql` file plus updates `atlas.sum`.

3. Review the generated `.sql`. It's goose-format:

   ```sql
   -- +goose Up
   ...

   -- +goose Down
   ...
   ```

4. Commit **both** the new `.sql` file and the updated `atlas.sum`.

5. Run the app — goose applies any pending migrations at startup:

   ```bash
   mise run dev
   ```

`atlas.sum` is Atlas's integrity checksum. goose ignores it, and it's excluded
from the `//go:embed *.sql` in `module.go`, so it never affects runtime.

## Hand-written migrations (data backfills, tricky changes)

Atlas covers schema diffs. When you need a **data migration** or a change Atlas
can't express, create an empty goose file and fill it in yourself:

```bash
mise run create-migration backfill_media_paths
```

Write the `-- +goose Up` / `-- +goose Down` blocks by hand. Keep it in the same
directory so goose applies it in timestamp order alongside the generated ones.

## Checking for unsafe changes (lint)

Before committing, lint the latest migration for destructive or non-portable
changes:

```bash
mise run db:lint     # atlas migrate lint --env bun --latest 1
```

## Down migrations are not auto-generated

Atlas does not produce versioned down-migrations, so the `-- +goose Down` block of
a **generated** file is a best-effort reverse and may be empty for some changes.
Devices roll **forward**, not back. If you genuinely need a rollback, hand-write
the `Down` block.

## CI drift check

CI verifies that the committed migrations are not behind the entities — i.e. that
nobody changed a struct without running `db:diff`:

```bash
mise run db:diff ci_drift_check
git diff --exit-code internal/database/migrations   # non-zero → entities are ahead
mise run db:lint
```

A non-zero `git diff` means an entity changed but no migration was generated for
it. Run `mise run db:diff <name>` locally, commit the result, and push.

## Troubleshooting

- **`db:diff` produces an empty/no migration** — the entities already match the
  migration history. Nothing to do.
- **`no such table` / loader finds nothing** — check `atlas.hcl`'s `--path` points
  at `./internal/database/entity` and every entity is in that package with a
  `bun.BaseModel` tag.
- **`checksum mismatch` from Atlas** — a `.sql` file was edited by hand after
  generation. Regenerate, or run `atlas migrate hash --env bun` to re-sync
  `atlas.sum` (only if you intentionally edited the file).
- **goose fails at startup** — a committed `.sql` has invalid SQL for SQLite.
  Fix the file (or the entity + regenerate) and restart.
