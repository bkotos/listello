# Adapter Repository Examples

Annotated references from the Listello codebase. Read when implementing a new or extended SQLite repository.

## 1. `SQLiteListRepository` — full implementation

**File:** `api/internal/adapter/list_repository.go`

Implements `application.ListRepository` with three methods: `Save`, `GetByID`, `GetAll`.

### Constructor

```go
type SQLiteListRepository struct {
	db *sql.DB
}

func NewSQLiteListRepository(sqlite *SQLite) *SQLiteListRepository {
	return &SQLiteListRepository{db: sqlite.db}
}
```

Takes `*SQLite` (not `*sql.DB` directly) — matches project wiring via `bootstrap.NewListService`.

### Save — upsert

```go
func (r *SQLiteListRepository) Save(list domain.List) error {
	const q = `
INSERT INTO lists (id, name) VALUES (?, ?)
ON CONFLICT(id) DO UPDATE SET name = excluded.name;`
	if _, err := r.db.Exec(q, list.ID, list.Name); err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	return nil
}
```

Key points:

- `ON CONFLICT DO UPDATE` for idempotent saves.
- Wrap DB errors with operation context (`save list`).

### GetByID — not found handling

```go
err := r.db.QueryRow(q, id).Scan(&listID, &name)
if err == sql.ErrNoRows {
	return domain.List{}, fmt.Errorf("list %q not found", id)
}
if err != nil {
	return domain.List{}, fmt.Errorf("find list: %w", err)
}
return domain.List{ID: listID, Name: name}, nil
```

Key points:

- `sql.ErrNoRows` becomes a readable not-found error (not wrapped with `%w`).
- Scan into locals, then construct domain type.

### GetAll — multi-row scan

```go
const q = `SELECT id, name FROM lists ORDER BY rowid`
// ... Query, defer rows.Close(), loop rows.Next(), check rows.Err()
```

Key points:

- `ORDER BY rowid` preserves insertion order (tested explicitly).
- Always `defer rows.Close()` and check `rows.Err()` after the loop.

## 2. Schema — `sqlite.go`

**File:** `api/internal/adapter/sqlite.go`

```go
func (s *SQLite) migrate() error {
	const q = `
CREATE TABLE IF NOT EXISTS lists (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL
);`
	if _, err := s.db.Exec(q); err != nil {
		return fmt.Errorf("migrate lists: %w", err)
	}
	return nil
}
```

`OpenSQLite` calls `migrate()` on open, so tests get schema automatically via `adapter.OpenSQLite`.

For new aggregates, add another `CREATE TABLE IF NOT EXISTS` block here. See [docs/sqlite-migrations.md](../../../docs/sqlite-migrations.md) when schema changes outgrow inline DDL.

## 3. Integration tests

**File:** `api/internal/adapter/list_repository_test.go`

Package: `adapter_test` (external test package).

### Save and read back

```go
func TestSQLiteListRepository_SaveAndGetByID(t *testing.T) {
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "lists.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	repo := adapter.NewSQLiteListRepository(db)
	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)

	require.NoError(t, repo.Save(list))
	got, err := repo.GetByID(list.ID)

	require.NoError(t, err)
	assert.Equal(t, list.ID, got.ID)
	assert.Equal(t, list.Name, got.Name)
}
```

Key points:

- Real SQLite in `t.TempDir()` — no mocks.
- Use domain helpers (`domain.CreateList`) to build valid test data.
- `t.Cleanup` closes the DB.

### GetAll with multiple rows

```go
func TestSQLiteListRepository_GetAll(t *testing.T) {
	// ... open db, create repo
	work, _, err := domain.CreateList("Work")
	personal, _, err := domain.CreateList("Personal")
	require.NoError(t, repo.Save(work))
	require.NoError(t, repo.Save(personal))

	got, err := repo.GetAll()

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, []domain.List{work, personal}, got)
}
```

Asserts order matches insertion order (`ORDER BY rowid`).

## 4. Application port (implemented by adapter)

**File:** `api/internal/application/list_service.go`

The adapter must match this interface:

```go
type ListRepository interface {
	Save(list domain.List) error
	GetAll() ([]domain.List, error)
	GetByID(id string) (domain.List, error)
}
```

`SQLiteListRepository` satisfies this implicitly — no `var _ application.ListRepository` compile-time check in the codebase today.

## 5. Bootstrap wiring (out of skill scope)

**File:** `api/internal/bootstrap/bootstrap.go`

```go
func NewListService(db *adapter.SQLite, eventLog *os.File) *application.ListService {
	lists := adapter.NewSQLiteListRepository(db)
	events := adapter.NewLoggingEventPublisher(eventLog)
	return application.NewListService(lists, events)
}
```

After implementing a new repository, mention that bootstrap needs a similar `New{Aggregate}Service` factory — do not implement unless asked.

## 6. Next repository: `ItemRepository`

**Application port** (`api/internal/application/item_service.go`):

```go
type ItemRepository interface {
	Save(listID string, item domain.Item) error
}
```

A `SQLiteItemRepository` would:

1. Add an `items` table in `sqlite.go` `migrate()`.
2. Implement `Save(listID string, item domain.Item) error` in `item_repository.go`.
3. Test with `adapter.OpenSQLite` + `domain` helpers to create items.

Use this as the template when the application port already exists but the adapter does not.
