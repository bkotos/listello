package sqlite_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bkotos/listello/internal/sqlite"
)

func TestParseEngine_SQLite(t *testing.T) {
	// Act
	got, err := sqlite.ParseEngine("sqlite")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, sqlite.EngineSQLite, got)
}

func TestParseEngine_EmptyDefaultsToSQLite(t *testing.T) {
	// Act
	got, err := sqlite.ParseEngine("")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, sqlite.EngineSQLite, got)
}

func TestParseEngine_Postgres(t *testing.T) {
	// Act
	got, err := sqlite.ParseEngine("postgres")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, sqlite.EnginePostgres, got)
}

func TestParseEngine_Unknown(t *testing.T) {
	// Act
	_, err := sqlite.ParseEngine("mysql")

	// Assert
	require.EqualError(t, err, `unsupported engine "mysql"`)
}
