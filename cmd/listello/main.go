package main

import (
	"fmt"
	"os"

	adapter "github.com/bkotos/listello/internal/listello-adapter"
	application "github.com/bkotos/listello/internal/listello-application"
)

func main() {
	name := "Next actions"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	db := openDB("listello.db")
	defer db.Close()

	eventLog := openEventLog("domain_events.log")
	defer eventLog.Close()

	listService := getListService(db, eventLog)

	listService.CreateList(name)
}

func getListService(db *adapter.SQLite, eventLog *os.File) *application.ListService {
	lists := adapter.NewSQLiteListRepository(db)
	events := adapter.NewLoggingEventPublisher(eventLog)
	listService := application.NewListService(lists, events)
	return listService
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
