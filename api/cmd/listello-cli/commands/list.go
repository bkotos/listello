package commands

import (
	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/listello-application"
)

func NewList(lists *application.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Manage lists",
	}
	cmd.AddCommand(NewListCreate(lists))
	return cmd
}
