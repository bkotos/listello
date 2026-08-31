package main

import (
	"fmt"
	"os"

	adapter "github.com/bkotos/listello/internal/listello-adapter"
	application "github.com/bkotos/listello/internal/listello-application"
)

func main() {
	db := openDB("listello.db")
	defer db.Close()

	eventLog := openEventLog("domain_events.log")
	defer eventLog.Close()

	listService := getListService(db, eventLog)
	if err := run(newRoot(listService)); err != nil {
		os.Exit(1)
	}
}

func getListService(db *adapter.SQLite, eventLog *os.File) *application.ListService {
	lists := adapter.NewSQLiteListRepository(db)
	events := adapter.NewLoggingEventPublisher(eventLog)
	return application.NewListService(lists, events)
}

func openDB(path string) *adapter.SQLite {
	db, err := adapter.OpenSQLite(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open sqlite: %v\n", err)
		os.Exit(1)
	}
	return db
}

func openEventLog(path string) *os.File {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open event log: %v\n", err)
		os.Exit(1)
	}
	return f
}
