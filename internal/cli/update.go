package cli

import (
	"fmt"

	"github.com/def4alt/markd/internal/bookmark"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(bookmarkSvc *bookmark.Service) *cobra.Command {
	var (
		title       string
		description string
		tags        []string
		url         string
	)

	updateCmd := &cobra.Command{
		Use:          "update",
		Short:        "Update a bookmark",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			b := &bookmark.UpdateInput{
				ID:          args[0],
				URL:         &url,
				Tags:        tags,
				Description: &description,
				Title:       &title,
			}

			err := bookmarkSvc.Update(cmd.Context(), b)
			if err != nil {
				return fmt.Errorf("update bookmark: %w", err)
			}

			return nil
		},
	}

	updateCmd.Flags().StringVarP(&title, "title", "t", "", "Bookmark title")
	updateCmd.Flags().StringVarP(&url, "url", "u", "", "Bookmark URL")
	updateCmd.Flags().StringVarP(&description, "description", "d", "", "Bookmark description")
	updateCmd.Flags().StringArrayVar(&tags, "tag", nil, "Bookmark tag")

	return updateCmd
}
