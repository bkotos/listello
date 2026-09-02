package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bkotos/listello/cmd/cli/commands"
	"github.com/bkotos/listello/internal/bootstrap"
)

func newRoot() (*cobra.Command, func()) {
	var dbPath string
	var cleanup func()

	c := &container{}
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&dbPath, "db", "listello.db", "SQLite database path")

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if c.list != nil {
			return nil
		}

		db := bootstrap.MustOpenDB(dbPath)
		eventLog := bootstrap.MustOpenEventLog("domain_events.log")
		c.list = bootstrap.NewListService(db, eventLog)
		c.item = bootstrap.NewItemService(db, eventLog)
		cleanup = func() {
			db.Close()
			eventLog.Close()
		}
		return nil
	}

	root.AddCommand(commands.NewList(c))
	root.AddCommand(commands.NewItem(c))
	return root, func() {
		if cleanup != nil {
			cleanup()
		}
	}
}

func run(root *cobra.Command) error {
	err := root.Execute()
	if err != nil {
		fmt.Fprintf(root.ErrOrStderr(), "error: %v\n", err)
	}
	return err
}
