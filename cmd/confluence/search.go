package confluence

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joselrodrigues/atlassian/internal/confluence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchCmd = &cobra.Command{
	Use:   "search [cql-query]",
	Short: "Search content using CQL",
	Long: `Search Confluence content using Confluence Query Language (CQL).

Examples:
  atlassian confluence search "space=MYSPACE"
  atlassian confluence search "type=page AND title~'Testing'"
  atlassian confluence search "text~'ABsmartly'" --limit 50`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cql := args[0]
		limit, _ := cmd.Flags().GetInt("limit")
		output := viper.GetString("output")

		client := confluence.NewClient()
		results, err := client.SearchContent(cql, []string{"space", "version"}, limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		printSearchResults(results, output)
	},
}

func init() {
	Cmd.AddCommand(searchCmd)
	searchCmd.Flags().Int("limit", 25, "Maximum number of results to return")
}

func printSearchResults(results *confluence.SearchResponse, format string) {
	if format == "json" {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("| %-12s | %-10s | %-50s |\n", "ID", "Space", "Title")
	fmt.Printf("| %-12s | %-10s | %-50s |\n", "------------", "----------", "--------------------------------------------------")
	for _, p := range results.Results {
		title := p.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}
		spaceKey := ""
		if p.Space != nil {
			spaceKey = p.Space.Key
		}
		fmt.Printf("| %-12s | %-10s | %-50s |\n", p.ID, spaceKey, title)
	}
	fmt.Printf("\nFound: %d results\n", results.Size)
}
