package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewItemTitle(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "title <item-id> <new-title>",
		Short: "Change an item's title",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]
			title := args[1]

			item, err := container.ItemService().ModifyItemTitle(itemID, title)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Renamed %s to %q\n", item.ID, item.Title)
			return nil
		},
	}
}
