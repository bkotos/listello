package bootstrap

import (
	"fmt"
	"os"

	adapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	"github.com/bkotos/listello/internal/sqlite"
)

// OpenDB opens a workspace database for the given engine and DSN.
func OpenDB(engine sqlite.Engine, dsn string) (*sqlite.WorkspaceDB, error) {
	workspace := sqlite.NewWorkspaceDB()
	if err := workspace.Open(engine, dsn); err != nil {
		return nil, err
	}
	return workspace, nil
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

// NewUserService wires user persistence and event publishing into UserService.
func NewUserService(workspace *sqlite.WorkspaceDB, eventLog *os.File) application.UserService {
	userRepo := adapter.NewSQLiteUserRepository(workspace)
	events := adapter.NewLoggingEventPublisher(eventLog)
	return application.NewUserService(userRepo, events)
}

// NewSpaceService wires space persistence and event publishing into SpaceService.
func NewSpaceService(workspace *sqlite.WorkspaceDB, eventLog *os.File) application.SpaceService {
	spaceRepo := adapter.NewSQLiteSpaceRepository(workspace)
	events := adapter.NewLoggingEventPublisher(eventLog)
	return application.NewSpaceService(spaceRepo, events)
}
