package sqlite

import (
	"database/sql"
	"fmt"
	"sync"
)

// WorkspaceDB holds the process-wide SQLite connection for workspace data.
// It starts closed and can be opened after persistence is initialized.
type WorkspaceDB struct {
	mu sync.RWMutex
	db *SQLite
}

// NewWorkspaceDB returns a closed workspace database holder.
func NewWorkspaceDB() *WorkspaceDB {
	return &WorkspaceDB{}
}

// Open opens (or reopens) the SQLite database at path.
func (w *WorkspaceDB) Open(path string) error {
	db, err := OpenSQLite(path)
	if err != nil {
		return fmt.Errorf("open workspace db: %w", err)
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if w.db != nil {
		_ = w.db.Close()
	}
	w.db = db
	return nil
}

// DB returns the underlying connection, or an error if persistence is not initialized.
func (w *WorkspaceDB) DB() (*sql.DB, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.db == nil {
		return nil, fmt.Errorf("persistence not initialized")
	}
	return w.db.DB(), nil
}

// Close closes the underlying database if open.
func (w *WorkspaceDB) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.db == nil {
		return nil
	}
	err := w.db.Close()
	w.db = nil
	return err
}
