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
			return fmt.Errorf("not implemented")
		},
	}
}
