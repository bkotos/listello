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

func TestRoot_RejectsPostgresUntilSupported(t *testing.T) {
	// Arrange
	root, cleanup := newRoot()
	defer cleanup()
	root.SetArgs([]string{"--engine", "postgres", "--db", "postgres://localhost/listello", "list", "create", "Test list"})

	// Act
	err := root.Execute()

	// Assert
	require.ErrorContains(t, err, `unsupported engine "postgres"`)
}
