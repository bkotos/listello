package commands

import (
	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/application"
)

func NewList(listService application.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Manage lists",
	}
	cmd.AddCommand(NewListCreate(listService))
	return cmd
}
