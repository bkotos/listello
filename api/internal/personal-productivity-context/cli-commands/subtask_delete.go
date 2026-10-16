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
			return fmt.Errorf("not implemented")
		},
	}
}
