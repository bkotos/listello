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
	return cmd
}
