package adapter_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bkotos/listello/internal/sqlite"
)

func openWorkspaceDB(t *testing.T, name string) *sqlite.WorkspaceDB {
	t.Helper()
	workspace := sqlite.NewWorkspaceDB()
	require.NoError(t, workspace.Open(filepath.Join(t.TempDir(), name)))
	t.Cleanup(func() { _ = workspace.Close() })
	return workspace
}
