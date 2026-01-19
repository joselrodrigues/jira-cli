package confluence

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joselrodrigues/atlassian/internal/confluence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "List pages in a space",
	Long:  `List all pages in a Confluence space.`,
	Run: func(cmd *cobra.Command, args []string) {
		spaceKey, _ := cmd.Flags().GetString("space")
		limit, _ := cmd.Flags().GetInt("limit")
		output := viper.GetString("output")

		if spaceKey == "" {
			fmt.Fprintln(os.Stderr, "Error: --space/-s flag is required")
			os.Exit(1)
		}

		client := confluence.NewClient()
		pages, err := client.GetSpaceContent(spaceKey, "page", limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		printPages(pages, output)
	},
}

func init() {
	Cmd.AddCommand(pagesCmd)
	pagesCmd.Flags().StringP("space", "s", "", "Space key (required)")
	pagesCmd.Flags().Int("limit", 25, "Maximum number of pages to return")
}

func printPages(pages *confluence.PageResults, format string) {
	if format == "json" {
		data, _ := json.MarshalIndent(pages, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("| %-12s | %-60s |\n", "ID", "Title")
	fmt.Printf("| %-12s | %-60s |\n", "------------", "------------------------------------------------------------")
	for _, p := range pages.Results {
		title := p.Title
		if len(title) > 60 {
			title = title[:57] + "..."
		}
		fmt.Printf("| %-12s | %-60s |\n", p.ID, title)
	}
	fmt.Printf("\nTotal: %d pages\n", pages.Size)
}
