package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSubtaskDelete(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <item-id> <subtask-id>",
		Short: "Delete a subtask on an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]
			subtaskID := args[1]

			item, err := container.ItemService().DeleteSubtask(itemID, subtaskID)
			if err != nil {
				return err
			}

			title := subtaskID
			for _, subtask := range item.Subtasks {
				if subtask.ID == subtaskID {
					title = subtask.Title
					break
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted subtask %q (%s) from %q\n", title, subtaskID, item.Title)
			return nil
		},
	}
}
