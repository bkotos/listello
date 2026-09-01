package main

import (
	"os"

	"github.com/bkotos/listello/internal/listello-bootstrap"
)

func main() {
	db := bootstrap.MustOpenDB("listello.db")
	defer db.Close()

	eventLog := bootstrap.MustOpenEventLog("domain_events.log")
	defer eventLog.Close()

	listService := bootstrap.NewListService(db, eventLog)
	if err := run(newRoot(listService)); err != nil {
		os.Exit(1)
	}
}
