# Adapter Repository Examples

Annotated references from the Listello codebase. Read when implementing a new or extended SQLite repository.

## 1. `SQLiteListRepository` — full implementation

**File:** `api/internal/personal-productivity-context/adapter/list_repository.go`

Implements `application.ListRepository` with three methods: `Save`, `GetByID`, `GetAll`.

### Constructor

```go
type SQLiteListRepository struct {
	workspace *sqlite.WorkspaceDB
}

func NewSQLiteListRepository(workspace *sqlite.WorkspaceDB) *SQLiteListRepository {
	return &SQLiteListRepository{workspace: workspace}
}
```

Takes `*sqlite.WorkspaceDB` — matches project wiring via `bootstrap.NewListService`. Queries go through bun (`openBun`).

### Save — upsert

```go
func (r *SQLiteListRepository) Save(list domain.List) error {
	db, err := openBun(r.workspace)
	if err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	_, err = db.NewInsert().
		Model(&listRow{ID: list.ID, Name: list.Name, CreatedAt: newCreatedAt()}).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	return nil
}
```

Key points:

- Bun `On("CONFLICT (id) DO UPDATE")` for idempotent saves.
- Omit `created_at` from the `Set` list so a later Save does not reshuffle GetAll order.
- Wrap DB errors with operation context (`save list`).

### GetByID — not found handling

```go
row := new(listRow)
err = db.NewSelect().Model(row).Where("id = ?", id).Scan(context.Background())
if errors.Is(err, sql.ErrNoRows) {
	return domain.List{}, fmt.Errorf("list %q not found", id)
}
if err != nil {
	return domain.List{}, fmt.Errorf("find list: %w", err)
}
return domain.List{ID: row.ID, Name: row.Name}, nil
```

Key points:

- `sql.ErrNoRows` becomes a readable not-found error (not wrapped with `%w`).
- Scan into a bun row struct, then construct the domain type.

### GetAll — multi-row scan

```go
var rows []listRow
err = db.NewSelect().Model(&rows).Order("created_at ASC", "id ASC").Scan(context.Background())
```

Key points:

- `Order("created_at ASC", "id ASC")` preserves insertion order (tested explicitly).
- Stamp `created_at` on insert only — omit it from the conflict `Set` list so a later Save does not reshuffle.

## 2. Schema — `sqlite.go`

**File:** `api/internal/personal-productivity-context/adapter/sqlite.go`

```go
func (s *SQLite) migrate() error {
	const q = `
CREATE TABLE IF NOT EXISTS lists (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	created_at TEXT NOT NULL
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

**File:** `api/internal/personal-productivity-context/adapter/list_repository_test.go`

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

Asserts order matches insertion order (`ORDER BY created_at, id`).

## 4. Application port (implemented by adapter)

**File:** `api/internal/personal-productivity-context/application/list_service.go`

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

**Application port** (`api/internal/personal-productivity-context/application/item_service.go`):

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
