package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListDelete(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <list-id>",
		Short: "Delete a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			listID := args[0]

			if err := container.ListService().DeleteList(listID); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Deleted list %s\n", listID)
			return nil
		},
	}
}
