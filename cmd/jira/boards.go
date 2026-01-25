package jira

import (
	"encoding/json"
	"fmt"

	"github.com/joselrodrigues/atlassian/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "List Jira boards",
	Long:  `List all Jira boards, optionally filtered by project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		output := viper.GetString("output")

		client := jira.NewClient()
		result, err := client.GetBoards(project)
		if err != nil {
			return fmt.Errorf("failed to get boards: %w", err)
		}

		if output == "json" {
			jsonData, err := json.MarshalIndent(result.Values, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(jsonData))
			return nil
		}

		if len(result.Values) == 0 {
			fmt.Println("No boards found")
			return nil
		}

		fmt.Println("| Board ID | Name | Type | Project |")
		fmt.Println("| -------- | ---- | ---- | ------- |")
		for _, board := range result.Values {
			projectKey := board.Location.ProjectKey
			if projectKey == "" {
				projectKey = "-"
			}
			fmt.Printf("| %d | %s | %s | %s |\n", board.ID, board.Name, board.Type, projectKey)
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(boardsCmd)

	boardsCmd.Flags().StringP("project", "p", "", "Filter by project key")
}
