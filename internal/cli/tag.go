package cli

import (
	"fmt"

	"github.com/def4alt/markd/internal/bookmark"
	"github.com/spf13/cobra"
)

func NewTagCmd(bookmarkSvc *bookmark.Service) *cobra.Command {
	tagCmd := &cobra.Command{
		Use:          "tag",
		Short:        "Tag a bookmark",
		Args:         cobra.MinimumNArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			b := &bookmark.UpdateInput{
				ID:   args[0],
				Tags: args[1:],
			}

			err := bookmarkSvc.Update(cmd.Context(), b)
			if err != nil {
				return fmt.Errorf("tag bookmark: %w", err)
			}

			return nil
		},
	}

	return tagCmd
}
