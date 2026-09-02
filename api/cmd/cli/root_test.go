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
