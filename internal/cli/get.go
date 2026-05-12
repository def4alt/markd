package cli

import (
	"fmt"

	"github.com/def4alt/markd/internal/bookmark"
	"github.com/spf13/cobra"
)

func NewGetCmd(bookmarkSvc *bookmark.Service) *cobra.Command {
	getCmd := &cobra.Command{
		Use:          "get",
		Short:        "Get a bookmark",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := bookmarkSvc.Get(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get bookmark: %w", err)
			}

			fmt.Printf("ID:          %s\n", b.ID)
			fmt.Printf("URL:         %s\n", b.URL)
			fmt.Printf("Title:       %s\n", b.Title)
			fmt.Printf("Description: %s\n", b.Description)
			fmt.Printf("Status:      %s\n", b.Status)
			fmt.Printf("Tags:        %v\n", b.Tags)
			fmt.Printf("Created:     %s\n", b.CreatedAt)
			fmt.Printf("Updated:     %s\n", b.UpdatedAt)

			return nil
		},
	}

	return getCmd
}
