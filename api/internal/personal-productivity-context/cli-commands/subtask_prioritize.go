package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func NewSubtaskPrioritize(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "prioritize <item-id> <subtask-id> <low|medium|high|none>",
		Short: "Change a subtask's priority on an item",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID := args[0]
			subtaskID := args[1]
			priority := subtaskPriorityArg(args[2])

			item, err := container.ItemService().PrioritizeSubtask(itemID, subtaskID, priority)
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
			fmt.Fprintf(cmd.OutOrStdout(), "Priority of subtask %q (%s) on %q changed to %s\n", title, subtaskID, item.Title, args[2])
			return nil
		},
	}
}

func subtaskPriorityArg(arg string) domain.ItemPriority {
	if arg == "none" {
		return domain.PriorityNone
	}
	return domain.ItemPriority(arg)
}
