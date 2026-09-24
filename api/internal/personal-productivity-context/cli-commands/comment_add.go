package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommentAdd(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "add <item-id> <text>",
		Short: "Comment on an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]
			body := args[1]

			instance, err := container.InstanceService().GetInstance()
			if err != nil {
				return err
			}

			item, err := container.ItemService().CommentItem(itemID, instance.User.ID, body)
			if err != nil {
				return err
			}

			commentID := item.Comments[len(item.Comments)-1].ID
			fmt.Fprintf(cmd.OutOrStdout(), "Commented on %q (%s)\n", item.Title, commentID)
			return nil
		},
	}
}
