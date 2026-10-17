package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommentDelete(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <item-id> <comment-id>",
		Short: "Delete a comment on an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]
			commentID := args[1]

			instance, err := container.InstanceService().GetInstance()
			if err != nil {
				return err
			}

			item, err := container.ItemService().DeleteComment(itemID, instance.User.ID, commentID)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Deleted comment %s from %q\n", commentID, item.Title)
			return nil
		},
	}
}
