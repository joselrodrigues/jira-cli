package jira

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joselrodrigues/atlassian/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update [issue-key]",
	Short: "Update an existing issue",
	Long: `Update fields of an existing Jira issue.
Supports updating summary, description, assignee, story points, and sprint.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueKey := args[0]
		summary, _ := cmd.Flags().GetString("summary")
		description, _ := cmd.Flags().GetString("description")
		fromStdin, _ := cmd.Flags().GetBool("stdin")
		assignee, _ := cmd.Flags().GetString("assignee")
		points, _ := cmd.Flags().GetFloat64("points")
		hasPoints := cmd.Flags().Changed("points")
		sprintID, _ := cmd.Flags().GetInt("sprint")
		storyPointsField, _ := cmd.Flags().GetString("points-field")

		if fromStdin {
			reader := bufio.NewReader(os.Stdin)
			var sb strings.Builder
			for {
				line, err := reader.ReadString('\n')
				sb.WriteString(line)
				if err != nil {
					break
				}
			}
			description = sb.String()
		}

		client := jira.NewClient()
		fields := make(map[string]interface{})

		if summary != "" {
			fields["summary"] = summary
		}
		if description != "" {
			fields["description"] = description
		}
		if hasPoints {
			fields[storyPointsField] = points
		}

		if assignee != "" {
			var accountID string
			if strings.Contains(assignee, "@") {
				users, err := client.SearchUsers(assignee)
				if err != nil {
					return fmt.Errorf("failed to search for user: %w", err)
				}
				if len(users) == 0 {
					return fmt.Errorf("no user found with email: %s", assignee)
				}
				accountID = users[0].AccountID
				fmt.Printf("Found user: %s (%s)\n", users[0].DisplayName, users[0].AccountID)
			} else {
				accountID = assignee
			}

			if err := client.AssignIssue(issueKey, accountID); err != nil {
				return fmt.Errorf("failed to assign issue: %w", err)
			}
			fmt.Printf("Issue %s assigned to %s\n", issueKey, accountID)
		}

		if sprintID != 0 {
			if err := client.MoveToSprint(sprintID, []string{issueKey}); err != nil {
				return fmt.Errorf("failed to move issue to sprint: %w", err)
			}
			fmt.Printf("Issue %s moved to sprint %d\n", issueKey, sprintID)
		}

		if len(fields) > 0 {
			if err := client.UpdateIssue(issueKey, fields); err != nil {
				return fmt.Errorf("failed to update issue: %w", err)
			}
		}

		if len(fields) == 0 && assignee == "" && sprintID == 0 {
			return fmt.Errorf("at least one field must be specified")
		}

		baseURL := strings.TrimSuffix(viper.GetString("jira_base_url"), "/")
		fmt.Printf("Issue %s updated successfully!\n", issueKey)
		fmt.Printf("URL: %s/browse/%s\n", baseURL, issueKey)

		return nil
	},
}

func init() {
	Cmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("summary", "s", "", "New summary")
	updateCmd.Flags().StringP("description", "d", "", "New description")
	updateCmd.Flags().Bool("stdin", false, "Read description from stdin")
	updateCmd.Flags().StringP("assignee", "a", "", "Assign to user (email or accountId)")
	updateCmd.Flags().Float64("points", 0, "Story points")
	updateCmd.Flags().Int("sprint", 0, "Sprint ID to move issue to")
	updateCmd.Flags().String("points-field", "customfield_10106", "Custom field ID for story points")
}
