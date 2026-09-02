package commands

import (
	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/application"
)

func NewItem(itemService application.ItemService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item",
		Short: "Manage items",
	}
	cmd.AddCommand(NewItemDefine(itemService))
	return cmd
}
