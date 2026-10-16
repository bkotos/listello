package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSubtaskUncomplete(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "uncomplete <item-id> <subtask-id>",
		Short: "Uncomplete a subtask on an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]
			subtaskID := args[1]

			item, err := container.ItemService().UncompleteSubtask(itemID, subtaskID)
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
			fmt.Fprintf(cmd.OutOrStdout(), "Uncompleted subtask %q (%s) on %q\n", title, subtaskID, item.Title)
			return nil
		},
	}
}
