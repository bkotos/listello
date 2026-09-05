package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewItemMove(container Container) *cobra.Command {
	var toListID string
	cmd := &cobra.Command{
		Use:   "move <item-id>",
		Short: "Move an item to another list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]

			item, err := container.ItemService().MoveItem(itemID, toListID)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Moved item %q (%s) to list %s\n", item.Title, item.ID, item.ListID)
			return nil
		},
	}
	cmd.Flags().StringVar(&toListID, "to", "", "destination list ID")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}
