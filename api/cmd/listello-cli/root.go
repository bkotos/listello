package main

import (
	"fmt"

	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/listello-application"
)

func newRoot(lists *application.ListService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newListCmd(lists))
	return root
}

func run(root *cobra.Command) error {
	err := root.Execute()
	if err != nil {
		fmt.Fprintf(root.ErrOrStderr(), "error: %v\n", err)
	}
	return err
}
