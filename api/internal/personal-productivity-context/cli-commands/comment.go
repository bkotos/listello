package commands

import (
	"github.com/spf13/cobra"
)

func NewComment(container Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Manage comments",
	}
	cmd.AddCommand(NewCommentAdd(container))
	return cmd
}
