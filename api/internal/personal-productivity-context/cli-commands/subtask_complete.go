package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSubtaskComplete(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "complete <item-id> <subtask-id>",
		Short: "Complete a subtask on an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]
			subtaskID := args[1]

			item, err := container.ItemService().CompleteSubtask(itemID, subtaskID)
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
			fmt.Fprintf(cmd.OutOrStdout(), "Completed subtask %q (%s) on %q\n", title, subtaskID, item.Title)
			return nil
		},
	}
}
