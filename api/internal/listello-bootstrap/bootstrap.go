package bootstrap

import (
	"fmt"
	"os"

	adapter "github.com/bkotos/listello/internal/listello-adapter"
	application "github.com/bkotos/listello/internal/listello-application"
)

// MustOpenDB opens SQLite or exits the process on failure.
func MustOpenDB(path string) *adapter.SQLite {
	db, err := adapter.OpenSQLite(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open sqlite: %v\n", err)
		os.Exit(1)
	}
	return db
}

// MustOpenEventLog opens the domain event log or exits the process on failure.
func MustOpenEventLog(path string) *os.File {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open event log: %v\n", err)
		os.Exit(1)
	}
	return f
}

// NewListService wires list persistence and event publishing into ListService.
func NewListService(db *adapter.SQLite, eventLog *os.File) *application.ListService {
	lists := adapter.NewSQLiteListRepository(db)
	events := adapter.NewLoggingEventPublisher(eventLog)
	return application.NewListService(lists, events)
}
