package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSubtaskLink(container Container) *cobra.Command {
	var parentID string
	cmd := &cobra.Command{
		Use:   "link <item-id>",
		Short: "Link an item as a child of another item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]

			item, err := container.ItemService().LinkAsChild(itemID, parentID)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Linked %q as a child of %s\n", item.Title, item.ParentID)
			return nil
		},
	}
	cmd.Flags().StringVar(&parentID, "parent", "", "parent item ID")
	_ = cmd.MarkFlagRequired("parent")
	return cmd
}
