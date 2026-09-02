package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bkotos/listello/cmd/cli/commands"
	application "github.com/bkotos/listello/internal/application"
)

func newRoot(listService application.ListService, itemService application.ItemService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(commands.NewList(listService))
	root.AddCommand(commands.NewItem(itemService))
	return root
}

func run(root *cobra.Command) error {
	err := root.Execute()
	if err != nil {
		fmt.Fprintf(root.ErrOrStderr(), "error: %v\n", err)
	}
	return err
}
