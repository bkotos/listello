package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewItemDelete(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <item-id>",
		Short: "Delete an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]

			if err := container.ItemService().DeleteItem(itemID); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Deleted item %s\n", itemID)
			return nil
		},
	}
}
