package commands

import (
	"github.com/spf13/cobra"
)

func NewSubtask(container Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subtask",
		Short: "Manage subtasks",
	}
	cmd.AddCommand(NewSubtaskAdd(container))
	cmd.AddCommand(NewSubtaskComplete(container))
	cmd.AddCommand(NewSubtaskUncomplete(container))
	cmd.AddCommand(NewSubtaskDelete(container))
	cmd.AddCommand(NewSubtaskPrioritize(container))
	return cmd
}
