package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoot_UsesCustomDBPath(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "custom.db")
	root, cleanup := newRoot()
	defer cleanup()
	stdout := &bytes.Buffer{}
	root.SetOut(stdout)
	root.SetArgs([]string{"--db", dbPath, "list", "create", "Test list"})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	_, err = os.Stat(dbPath)
	require.NoError(t, err)
}

func TestRoot_RejectsUnknownEngine(t *testing.T) {
	// Arrange
	root, cleanup := newRoot()
	defer cleanup()
	root.SetArgs([]string{"--engine", "mysql", "list", "create", "Test list"})

	// Act
	err := root.Execute()

	// Assert
	require.EqualError(t, err, `unsupported engine "mysql"`)
}

func TestRoot_PostgresEngineIsNoLongerRejected(t *testing.T) {
	// Arrange
	root, cleanup := newRoot()
	defer cleanup()
	root.SetArgs([]string{"--engine", "postgres", "--db", "postgres://listello:listello@127.0.0.1:1/listello?sslmode=disable", "list", "create", "Test list"})

	// Act
	err := root.Execute()

	// Assert
	require.Error(t, err)
	require.NotContains(t, err.Error(), `unsupported engine "postgres"`)
}
