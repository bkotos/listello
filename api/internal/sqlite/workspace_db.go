package sqlite

import (
	"database/sql"
	"fmt"
	"sync"
)

// WorkspaceDB holds the process-wide workspace database connection.
// It starts closed and can be opened after persistence is initialized.
type WorkspaceDB struct {
	mu     sync.RWMutex
	engine Engine
	db     *SQLite
}

// NewWorkspaceDB returns a closed workspace database holder.
func NewWorkspaceDB() *WorkspaceDB {
	return &WorkspaceDB{}
}

// Open opens (or reopens) the workspace database for the given engine and DSN.
// For SQLite, dsn is a file path. Postgres is recognized but not opened yet.
func (w *WorkspaceDB) Open(engine Engine, dsn string) error {
	switch engine {
	case EngineSQLite:
	default:
		return fmt.Errorf("unsupported engine %q", engine)
	}

	db, err := OpenSQLite(dsn)
	if err != nil {
		return fmt.Errorf("open workspace db: %w", err)
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if w.db != nil {
		_ = w.db.Close()
	}
	w.engine = engine
	w.db = db
	return nil
}

// Engine returns the open engine, or an error if persistence is not initialized.
func (w *WorkspaceDB) Engine() (Engine, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.db == nil {
		return "", fmt.Errorf("persistence not initialized")
	}
	return w.engine, nil
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
	w.engine = ""
	return err
}
