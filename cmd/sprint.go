package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/joserodrigues/jira-cli/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sprintCmd = &cobra.Command{
	Use:   "sprint",
	Short: "List issues in current sprint",
	Long:  `List all issues in the current active sprint for a project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		client := jira.NewClient()
		result, err := client.GetSprintIssues(project)
		if err != nil {
			return fmt.Errorf("failed to get sprint issues: %w", err)
		}

		if viper.GetString("output") == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Sprint issues for project %s:\n\n", project)
		printSearchResults(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sprintCmd)
	sprintCmd.Flags().StringP("project", "p", "", "Project key (required)")
	sprintCmd.MarkFlagRequired("project")
}