package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/joselrodrigues/jira-cli/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchCmd = &cobra.Command{
	Use:   "search [jql]",
	Short: "Search issues using JQL",
	Long:  `Search for Jira issues using JQL (Jira Query Language).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jql := args[0]
		maxResults, _ := cmd.Flags().GetInt("max")

		client := jira.NewClient()
		result, err := client.SearchIssues(jql, maxResults)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if viper.GetString("output") == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		printSearchResults(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntP("max", "m", 50, "Maximum results to return")
}

func printSearchResults(result *jira.SearchResult) {
	fmt.Printf("Found %d issues:\n\n", result.Total)
	fmt.Printf("| Key | Status | SP | Assignee | Summary |\n")
	fmt.Printf("|-----|--------|-----|----------|--------|\n")

	for _, issue := range result.Issues {
		assignee := "Unassigned"
		if issue.Fields.Assignee != nil {
			assignee = issue.Fields.Assignee.DisplayName
		}

		sp := "-"
		if issue.Fields.StoryPoints > 0 {
			sp = fmt.Sprintf("%.0f", issue.Fields.StoryPoints)
		}

		summary := issue.Fields.Summary
		if len(summary) > 50 {
			summary = summary[:47] + "..."
		}

		fmt.Printf("| %s | %s | %s | %s | %s |\n",
			issue.Key,
			issue.Fields.Status.Name,
			sp,
			assignee,
			summary,
		)
	}
}