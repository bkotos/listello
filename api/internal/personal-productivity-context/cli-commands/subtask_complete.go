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
			return fmt.Errorf("not implemented")
		},
	}
}
