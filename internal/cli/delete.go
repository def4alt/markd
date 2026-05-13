package cli

import (
	"fmt"

	"github.com/def4alt/markd/internal/bookmark"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(bookmarkSvc *bookmark.Service) *cobra.Command {
	getCmd := &cobra.Command{
		Use:          "delete",
		Short:        "Delete a bookmark",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := bookmarkSvc.Delete(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("delete bookmark: %w", err)
			}

			fmt.Printf("Successfully removed %s", args[0])

			return nil
		},
	}

	return getCmd
}
