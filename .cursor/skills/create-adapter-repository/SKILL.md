---
name: create-adapter-repository
description: >-
  Scaffolds and implements SQLite repository adapters in api/internal/adapter/
  that satisfy application-layer ports. Use when adding or extending adapter
  repositories, SQLite persistence, schema migrations, or SQLiteListRepository/
  SQLiteItemRepository implementations following Listello hexagonal architecture.
---

# Create Adapter Repository

Guide for implementing repository adapters in `api/internal/adapter/`. Read this skill before changing adapter persistence code.

Architecture context: see [README.md](../../../README.md) (adapters implement application ports and are the only layer that talks to SQLite). Layer order: [LAYER-ORDER.md](../LAYER-ORDER.md). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc).

## Upstream dependencies

**Stop and do not proceed** until application-layer code exists. See [LAYER-ORDER.md](../LAYER-ORDER.md).

Before starting, verify:

| Check | How |
|-------|-----|
| `{Aggregate}Repository` interface exists | `internal/application/{aggregate}_service.go` |
| `{Aggregate}Service` + constructor exist | Same file |
| Interface includes methods you will implement | `Save`, `GetByID`, etc. |

If any check fails → **stop**. Tell the user to use `create-application-service` first. Do not write adapter tests or SQL.

## Scope

**In scope:** `SQLite{Aggregate}Repository`, SQL queries, schema additions in `sqlite.go`, integration tests with real SQLite.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- Application port interfaces (`api/internal/application/`) — use [create-application-service](../create-application-service/SKILL.md) first
- Bootstrap wiring (`api/internal/bootstrap/bootstrap.go`)
- HTTP handlers (`api/cmd/server/handlers/`)
- CLI commands (`api/cmd/cli/commands/`)

## Architecture constraints

- Adapter implements an application port interface (e.g. `application.ListRepository`).
- Adapter depends on `internal/domain` and `database/sql` — never on HTTP, CLI, or application service types.
- Map SQL rows to domain types; no business rules in the adapter.
- Wrap errors with context: `fmt.Errorf("save list: %w", err)`.
- Return domain-friendly errors for not-found (`sql.ErrNoRows` → e.g. `fmt.Errorf("list %q not found", id)`).

## Preconditions

These duplicate the upstream gate — all must pass:

1. **Application port exists** — interface in `api/internal/application/{aggregate}_service.go` with the methods to implement.
2. **Domain types defined** — know which fields to persist and scan.
3. **Schema planned** — new table or column changes identified.

If the port is missing, **stop** and use `create-application-service` first.

## Decision tree

1. **New aggregate repository?** → Create `{aggregate}_repository.go` + `{aggregate}_repository_test.go`; add table to `sqlite.go` `migrate()`.
2. **Existing repository, new method?** → Add method to existing adapter file; add test; extend schema if needed.
3. **Complex schema change?** → See [docs/sqlite-migrations.md](../../../docs/sqlite-migrations.md) (goose). For now, `CREATE TABLE IF NOT EXISTS` in `migrate()` is the project default.

## Scaffold checklist

```
Task progress:
- [ ] Verify upstream: application port + service exist (stop if not — see LAYER-ORDER.md)
- [ ] Read application port interface (method signatures)
- [ ] Read domain type fields to persist
- [ ] Decide: new repository vs extend existing
- [ ] Add CREATE TABLE (or ALTER) to sqlite.go migrate() if needed
- [ ] Add SQLite{Aggregate}Repository struct + constructor
- [ ] Add method stubs returning fmt.Errorf("not implemented") for red phase
- [ ] Write failing integration test(s) in adapter_test package
- [ ] Run tests — confirm failure is missing behavior
- [ ] STOP — summarize spec and failure for user review
- [ ] (After approval) Implement minimum production code
- [ ] Re-run tests — confirm green
```

## Naming conventions

| Artifact | Convention |
|----------|------------|
| Repository file | `{aggregate}_repository.go` |
| Test file | `{aggregate}_repository_test.go` |
| Test package | `adapter_test` |
| Implementation type | `SQLite{Aggregate}Repository` |
| Constructor | `NewSQLite{Aggregate}Repository(sqlite *SQLite)` |
| Test name | `TestSQLite{Aggregate}Repository_{Method}_{Behavior}` |
| Domain import | `domain "github.com/bkotos/listello/internal/domain"` |

Existing shared infrastructure: `SQLite`, `OpenSQLite` in `sqlite.go`.

## Code templates

### New repository struct

```go
package adapter

import (
	"database/sql"
	"fmt"

	domain "github.com/bkotos/listello/internal/domain"
)

// SQLite{Aggregate}Repository persists {aggregates} in SQLite.
type SQLite{Aggregate}Repository struct {
	db *sql.DB
}

// NewSQLite{Aggregate}Repository returns a {aggregate} repository using the given SQLite connection.
func NewSQLite{Aggregate}Repository(sqlite *SQLite) *SQLite{Aggregate}Repository {
	return &SQLite{Aggregate}Repository{db: sqlite.db}
}
```

### Save (upsert)

```go
// Save stores the {aggregate}.
func (r *SQLite{Aggregate}Repository) Save(/* args */) error {
	const q = `
INSERT INTO {table} (/* columns */) VALUES (/* ? placeholders */)
ON CONFLICT(/* pk */) DO UPDATE SET /* columns = excluded.columns */;`
	if _, err := r.db.Exec(q, /* values */); err != nil {
		return fmt.Errorf("save {aggregate}: %w", err)
	}
	return nil
}
```

### GetByID

```go
// GetByID returns the {aggregate} with the given ID.
func (r *SQLite{Aggregate}Repository) GetByID(id string) (domain.{Aggregate}, error) {
	const q = `SELECT /* columns */ FROM {table} WHERE id = ?`
	// scan into locals
	err := r.db.QueryRow(q, id).Scan(/* &fields */)
	if err == sql.ErrNoRows {
		return domain.{Aggregate}{}, fmt.Errorf("{aggregate} %q not found", id)
	}
	if err != nil {
		return domain.{Aggregate}{}, fmt.Errorf("find {aggregate}: %w", err)
	}
	return domain.{Aggregate}{/* fields */}, nil
}
```

### GetAll (multi-row)

```go
// GetAll returns all {aggregates} in insertion order.
func (r *SQLite{Aggregate}Repository) GetAll() ([]domain.{Aggregate}, error) {
	const q = `SELECT /* columns */ FROM {table} ORDER BY rowid`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("list {aggregates}: %w", err)
	}
	defer rows.Close()

	var results []domain.{Aggregate}
	for rows.Next() {
		// scan row, append to results
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list {aggregates}: %w", err)
	}
	return results, nil
}
```

### Schema in sqlite.go

Add to `migrate()`:

```go
const q = `
CREATE TABLE IF NOT EXISTS {table} (
	/* columns */
);`
if _, err := s.db.Exec(q); err != nil {
	return fmt.Errorf("migrate {table}: %w", err)
}
```

### Red-phase stub

```go
return fmt.Errorf("not implemented")
// or for methods returning a value:
return domain.{Aggregate}{}, fmt.Errorf("not implemented")
```

## Test templates

Use `adapter_test` package with real SQLite — no mocks. Open a temp database per test.

### Test setup

```go
func TestSQLite{Aggregate}Repository_{Method}_{Behavior}(t *testing.T) {
	// Arrange
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	repo := adapter.NewSQLite{Aggregate}Repository(db)
	// seed data via domain helpers or repo.Save

	// Act
	// ...

	// Assert
}
```

### Save and read back

```go
func TestSQLiteListRepository_SaveAndGetByID(t *testing.T) {
	// Arrange
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "lists.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	repo := adapter.NewSQLiteListRepository(db)
	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)

	// Act
	require.NoError(t, repo.Save(list))
	got, err := repo.GetByID(list.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, list.ID, got.ID)
	assert.Equal(t, list.Name, got.Name)
}
```

One behavioral concern per test. Prefer round-trip tests (save → read) over testing SQL in isolation.

## TDD workflow

Follow `.cursor/rules/tdd.mdc`:

1. **Red** — Write failing integration test first. Production may be a `not implemented` stub. Run tests; confirm failure is missing behavior, not a broken harness.
2. **Stop** — Summarize the new/changed spec and the failure. **Do not implement** until the user explicitly approves.
3. **Green** — Implement the bare minimum to pass. No extra queries, indexes, or error cases the test does not assert. Re-run tests.

Do not write production implementation in the same turn as a new failing spec.

## Schema migrations

- **Simple (default):** add `CREATE TABLE IF NOT EXISTS` to `sqlite.go` `migrate()`.
- **Complex changes** (rename column, backfill, drop): see [docs/sqlite-migrations.md](../../../docs/sqlite-migrations.md) for goose guidance.

Do not break existing tables without a migration plan.

## Verification

```bash
# All API tests
make -C api test

# Adapter layer only
cd api && go test ./internal/adapter/...
```

Confirm the implementation satisfies the application port by compiling packages that wire them (bootstrap can be updated separately).

## Further reading

Annotated walkthrough of `SQLiteListRepository`: [examples.md](examples.md)
