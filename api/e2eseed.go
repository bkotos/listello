//go:build ignore

package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	instanceadapter "github.com/bkotos/listello/internal/listello-instance-context/adapter"
	"github.com/bkotos/listello/internal/sqlite"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run e2eseed.go <db-path>\n")
		os.Exit(1)
	}

	const (
		userID   = "US_e2e"
		userName = "Alex"
	)

	db, err := sqlite.OpenSQLite(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "open sqlite: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	if _, err := db.DB().Exec(`INSERT INTO users (id, name) VALUES (?, ?)`, userID, userName); err != nil {
		fmt.Fprintf(os.Stderr, "insert user: %v\n", err)
		os.Exit(1)
	}

	locatorPath, err := instanceadapter.ListelloInstanceLocatorPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "locator path: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(locatorPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir locator: %v\n", err)
		os.Exit(1)
	}
	f, err := os.Create(locatorPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create locator: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	err = gob.NewEncoder(f).Encode(instanceadapter.ListelloInstanceFile{
		SchemaVersion: instanceadapter.ListelloInstanceSchemaVersion,
		Data: instanceadapter.ListelloInstanceData{
			ID: "LI_e2e",
			User: instanceadapter.UserData{
				ID:   userID,
				Name: userName,
			},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode instance: %v\n", err)
		os.Exit(1)
	}
}
