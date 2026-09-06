package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListCreate(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := container.ListService().CreateList(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created list %q (%s)\n", list.Name, list.ID)
			return nil
		},
	}
}
