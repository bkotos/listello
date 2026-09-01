package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd(lists listService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Manage lists",
	}
	cmd.AddCommand(newListCreateCmd(lists))
	return cmd
}

func newListCreateCmd(lists listService) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := lists.CreateList(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created list %q (%s)\n", list.Name, list.ID)
			return nil
		},
	}
}
