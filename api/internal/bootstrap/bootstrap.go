package bootstrap

import (
	"fmt"
	"os"

	adapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	"github.com/bkotos/listello/internal/sqlite"
)

// MustOpenDB opens a workspace SQLite database or exits the process on failure.
func MustOpenDB(path string) *sqlite.WorkspaceDB {
	workspace := sqlite.NewWorkspaceDB()
	if err := workspace.Open(path); err != nil {
		fmt.Fprintf(os.Stderr, "open sqlite: %v\n", err)
		os.Exit(1)
	}
	return workspace
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
func NewListService(workspace *sqlite.WorkspaceDB, eventLog *os.File) application.ListService {
	lists := adapter.NewSQLiteListRepository(workspace)
	events := adapter.NewLoggingEventPublisher(eventLog)
	return application.NewListService(lists, events)
}

// NewItemService wires list and item persistence and event publishing into ItemService.
func NewItemService(workspace *sqlite.WorkspaceDB, eventLog *os.File) application.ItemService {
	listRepo := adapter.NewSQLiteListRepository(workspace)
	itemRepo := adapter.NewSQLiteItemRepository(workspace)
	events := adapter.NewLoggingEventPublisher(eventLog)
	return application.NewItemService(listRepo, itemRepo, events)
}
