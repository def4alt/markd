package cli

import (
	"github.com/def4alt/markd/internal/bookmark"
	"github.com/spf13/cobra"
)

func NewRootCmd(bookmarkSvc *bookmark.Service) *cobra.Command {
	root := &cobra.Command{
		Use:   "markctl",
		Short: "Manage bookmarks",
	}

	root.AddCommand(NewAddCmd(bookmarkSvc))
	root.AddCommand(NewGetCmd(bookmarkSvc))
	root.AddCommand(NewListCmd(bookmarkSvc))
	root.AddCommand(NewTagCmd(bookmarkSvc))
	root.AddCommand(NewSearchCmd(bookmarkSvc))
	root.AddCommand(NewDeleteCmd(bookmarkSvc))

	return root
}
