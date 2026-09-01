package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/listello-application"
)

func NewListCreate(lists *application.ListService) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := lists.CreateList(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created list %q (%s)\n", list.Name, list.ID)
			return nil
		},
	}
}
