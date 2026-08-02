# SQLite schema migrations

Notes for evolving the Listello SQLite schema (new tables, column changes, etc.).

## Current state

`adapter.OpenSQLite` runs a simple `migrate()` with `CREATE TABLE IF NOT EXISTS`. That is fine early on for:

- local `listello.db`
- additive “create new table” changes
- throwing away the DB when iterating

It does **not** handle well:

- renaming/dropping columns
- changing constraints/types
- one-time data backfills

## Recommended approach: Goose as a library

Use [pressly/goose](https://github.com/pressly/goose) **from Go code** (not only the CLI) so migrations run when the app starts (or via a dedicated command). That way a published binary can migrate the user’s DB without requiring goose to be installed.

### Pattern

1. Keep numbered SQL migrations, e.g. `migrations/00001_lists.sql`, `00002_items.sql`.
2. Embed them with `embed.FS`.
3. On open / boot: open SQLite → run goose against the embedded FS → then use repositories.
4. Goose records applied versions (typically a `goose_db_version` table) so each migration runs once, in order.

### Why this fits Listello

- Schema can evolve safely as aggregates grow (items, etc.).
- Same path works for `make run` and for a released binary.
- CLI goose remains optional for local/dev convenience.

### Alternatives (not preferred long-term)

- Keep expanding a hand-rolled `migrate()` with `CREATE TABLE IF NOT EXISTS` / careful `ALTER TABLE`s while OK deleting `listello.db`.
- ORM “automigrate” — usually a poor fit for explicit DDD-style schema history.

## When to do this

Introduce versioned goose migrations when adding the next table or changing `lists` meaningfully. Move the current `lists` DDL into the first migration file and replace the inline `migrate()` in `sqlite.go` with a goose library call.
