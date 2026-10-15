package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSubtaskAdd(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "add <item-id> <title>",
		Short: "Add a subtask to an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]
			title := args[1]

			item, err := container.ItemService().AddSubtask(itemID, title)
			if err != nil {
				return err
			}

			subtaskID := item.Subtasks[len(item.Subtasks)-1].ID
			fmt.Fprintf(cmd.OutOrStdout(), "Added subtask %q (%s) to %q\n", title, subtaskID, item.Title)
			return nil
		},
	}
}
