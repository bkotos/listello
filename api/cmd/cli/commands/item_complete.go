package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewItemComplete(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "complete <item-id>",
		Short: "Complete an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]

			item, err := container.ItemService().CompleteItem(itemID)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Completed item %q (%s)\n", item.Title, item.ID)
			return nil
		},
	}
}
