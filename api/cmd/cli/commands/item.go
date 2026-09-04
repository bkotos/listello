package commands

import (
	"github.com/spf13/cobra"
)

func NewItem(container Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item",
		Short: "Manage items",
	}
	cmd.AddCommand(NewItemDefine(container))
	cmd.AddCommand(NewItemComplete(container))
	cmd.AddCommand(NewItemDelete(container))
	cmd.AddCommand(NewItemTitle(container))
	cmd.AddCommand(NewItemList(container))
	return cmd
}
