package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/joselrodrigues/jira-cli/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get [issue-key]",
	Short: "Get issue details",
	Long:  `Retrieve detailed information about a Jira issue by its key (e.g., PROJECT-123).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueKey := args[0]
		client := jira.NewClient()

		issue, err := client.GetIssue(issueKey)
		if err != nil {
			return fmt.Errorf("failed to get issue: %w", err)
		}

		if viper.GetString("output") == "json" {
			data, _ := json.MarshalIndent(issue, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		printIssue(issue)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func printIssue(issue *jira.Issue) {
	fmt.Printf("## %s\n\n", issue.Key)
	fmt.Printf("| Campo | Valor |\n")
	fmt.Printf("|-------|-------|\n")
	fmt.Printf("| **Summary** | %s |\n", issue.Fields.Summary)
	fmt.Printf("| **Status** | %s |\n", issue.Fields.Status.Name)
	fmt.Printf("| **Priority** | %s |\n", issue.Fields.Priority.Name)

	if issue.Fields.Assignee != nil {
		fmt.Printf("| **Assignee** | %s |\n", issue.Fields.Assignee.DisplayName)
	} else {
		fmt.Printf("| **Assignee** | Unassigned |\n")
	}

	if issue.Fields.StoryPoints > 0 {
		fmt.Printf("| **Story Points** | %.0f |\n", issue.Fields.StoryPoints)
	}

	if issue.Fields.Description != "" {
		fmt.Printf("\n### Description\n\n%s\n", issue.Fields.Description)
	}
}