package commands

import (
	"github.com/spf13/cobra"
)

func NewList(container Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Manage lists",
	}
	cmd.AddCommand(NewListCreate(container))
	return cmd
}
