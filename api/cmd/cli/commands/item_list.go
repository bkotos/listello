package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/application"
	domain "github.com/bkotos/listello/internal/domain"
)

func NewItemList(itemService application.ItemService) *cobra.Command {
	return &cobra.Command{
		Use:   "list <list-id>",
		Short: "List items on a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			listID := args[0]

			items, err := itemService.GetAll(listID)
			if err != nil {
				return err
			}

			for _, item := range items {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s  %s\n", itemCheckbox(item.State), item.ID, item.Title)
			}
			return nil
		},
	}
}

func itemCheckbox(state domain.ItemState) string {
	if state == domain.ItemComplete {
		return "[x]"
	}
	return "[ ]"
}
