package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/application"
)

func NewItemDefine(itemService application.ItemService) *cobra.Command {
	return &cobra.Command{
		Use:   "define <list-id> <title>",
		Short: "Define an item on a list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			listID := args[0]
			title := args[1]

			item, err := itemService.DefineItem(listID, title)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Defined item %q (%s) on list %s\n", item.Title, item.ID, listID)
			return nil
		},
	}
}
