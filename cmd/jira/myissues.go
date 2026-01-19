package jira

import (
	"encoding/json"
	"fmt"

	"github.com/joselrodrigues/atlassian/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var myIssuesCmd = &cobra.Command{
	Use:   "my-issues",
	Short: "List my assigned open issues",
	Long:  `List all open issues assigned to the current user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := jira.NewClient()
		result, err := client.GetMyIssues()
		if err != nil {
			return fmt.Errorf("failed to get issues: %w", err)
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
	Cmd.AddCommand(myIssuesCmd)
}
