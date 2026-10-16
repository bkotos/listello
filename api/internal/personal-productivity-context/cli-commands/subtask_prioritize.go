package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSubtaskPrioritize(container Container) *cobra.Command {
	return &cobra.Command{
		Use:   "prioritize <item-id> <subtask-id> <low|medium|high|none>",
		Short: "Change a subtask's priority on an item",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}
}
