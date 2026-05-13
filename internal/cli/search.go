package cli

import (
	"fmt"

	"github.com/def4alt/markd/internal/bookmark"
	"github.com/spf13/cobra"
)

func NewSearchCmd(bookmarkSvc *bookmark.Service) *cobra.Command {
	tagCmd := &cobra.Command{
		Use:          "search",
		Short:        "Search a bookmark",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := bookmarkSvc.Search(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("search bookmark: %w", err)
			}

			for _, bm := range b {
				fmt.Printf("%s: %s\n", bm.ID, bm.URL)
			}

			return nil
		},
	}

	return tagCmd
}
