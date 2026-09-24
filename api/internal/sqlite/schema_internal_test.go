package sqlite

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

type stubResult struct{}

func (stubResult) LastInsertId() (int64, error) { return 0, nil }
func (stubResult) RowsAffected() (int64, error) { return 0, nil }

type catalogRaceExecer struct {
	failsLeft int
	calls     int
}

func (e *catalogRaceExecer) Exec(query string, args ...any) (sql.Result, error) {
	e.calls++
	if e.failsLeft > 0 {
		e.failsLeft--
		return nil, &pgconn.PgError{
			Code:    "23505",
			Message: `duplicate key value violates unique constraint "pg_type_typname_nsp_index"`,
		}
	}
	return stubResult{}, nil
}

type alwaysErrorExecer struct {
	err error
}

func (e alwaysErrorExecer) Exec(query string, args ...any) (sql.Result, error) {
	return nil, e.err
}

func TestMigrateSchema_RetriesPostgresCatalogUniqueViolation(t *testing.T) {
	db := &catalogRaceExecer{failsLeft: 1}

	require.NoError(t, migrateSchema(db))
	require.Greater(t, db.calls, 1)
}

func TestMigrateSchema_ReturnsNonCatalogErrors(t *testing.T) {
	want := errors.New("connection reset")

	err := migrateSchema(alwaysErrorExecer{err: want})

	require.EqualError(t, err, "migrate lists: connection reset")
}
