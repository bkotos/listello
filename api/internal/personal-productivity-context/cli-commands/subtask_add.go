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
			return fmt.Errorf("not implemented")
		},
	}
}
