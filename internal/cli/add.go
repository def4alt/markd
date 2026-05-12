package cli

import (
	"fmt"

	"github.com/def4alt/markd/internal/bookmark"
	"github.com/spf13/cobra"
)

func NewAddCmd(bookmarkSvc *bookmark.Service) *cobra.Command {
	var (
		title       string
		description string
		tags        []string
	)

	addCmd := &cobra.Command{
		Use:          "add",
		Short:        "Add a bookmark",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := bookmarkSvc.Add(cmd.Context(), bookmark.AddInput{
				URL:         args[0],
				Title:       title,
				Description: description,
				Tags:        tags,
			})
			if err != nil {
				return fmt.Errorf("add bookmark: %w", err)
			}

			return nil
		},
	}

	addCmd.Flags().StringVarP(&title, "title", "t", "", "Bookmark title")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Bookmark description")
	addCmd.Flags().StringArrayVar(&tags, "tag", nil, "Bookmark tag")

	return addCmd
}
