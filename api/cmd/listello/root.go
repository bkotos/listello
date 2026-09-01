package main

import (
	"fmt"

	"github.com/spf13/cobra"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

type listService interface {
	CreateList(name string) (domain.List, error)
}

func newRoot(lists listService) *cobra.Command {
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
