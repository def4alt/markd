package cli

import (
	"fmt"

	"github.com/def4alt/markd/internal/bookmark"
	"github.com/spf13/cobra"
)

func NewListCmd(bookmarkSvc *bookmark.Service) *cobra.Command {
	listCmd := &cobra.Command{
		Use:          "list",
		Short:        "List bookmarks",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			bookmarks, err := bookmarkSvc.List(cmd.Context())
			if err != nil {
				return fmt.Errorf("list bookmark: %w", err)
			}

			for _, b := range bookmarks {
				fmt.Printf("%s  %s  %s  %s\n", b.ID, b.URL, b.Title, b.Status)
			}

			return nil
		},
	}

	return listCmd
}
