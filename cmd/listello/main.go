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

	db, err := adapter.OpenSQLite("listello.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "open sqlite: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	eventLog, err := os.OpenFile("domain_events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open event log: %v\n", err)
		os.Exit(1)
	}
	defer eventLog.Close()

	lists := adapter.NewSQLiteListRepository(db)
	events := adapter.NewLoggingEventPublisher(eventLog)
	svc := application.NewListService(lists, events)

	list, err := svc.CreateList(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create list: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("created list %q (%s)\n", list.Name, list.ID)
}
