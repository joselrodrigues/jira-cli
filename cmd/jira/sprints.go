package jira

import (
	"encoding/json"
	"fmt"

	"github.com/joselrodrigues/atlassian/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sprintsCmd = &cobra.Command{
	Use:   "sprints",
	Short: "List sprints for a board",
	Long:  `List all sprints for a specific board, optionally filtered by state.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetInt("board")
		state, _ := cmd.Flags().GetString("state")
		output := viper.GetString("output")

		if boardID == 0 {
			return fmt.Errorf("--board is required")
		}

		client := jira.NewClient()
		result, err := client.GetSprints(boardID, state)
		if err != nil {
			return fmt.Errorf("failed to get sprints: %w", err)
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
			fmt.Println("No sprints found")
			return nil
		}

		fmt.Println("| Sprint ID | Name | State | Start Date | End Date |")
		fmt.Println("| --------- | ---- | ----- | ---------- | -------- |")
		for _, sprint := range result.Values {
			startDate := sprint.StartDate
			if startDate == "" {
				startDate = "-"
			} else if len(startDate) > 10 {
				startDate = startDate[:10]
			}
			endDate := sprint.EndDate
			if endDate == "" {
				endDate = "-"
			} else if len(endDate) > 10 {
				endDate = endDate[:10]
			}
			fmt.Printf("| %d | %s | %s | %s | %s |\n", sprint.ID, sprint.Name, sprint.State, startDate, endDate)
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(sprintsCmd)

	sprintsCmd.Flags().IntP("board", "b", 0, "Board ID (required)")
	sprintsCmd.Flags().StringP("state", "s", "", "Filter by state (active, future, closed)")
}
