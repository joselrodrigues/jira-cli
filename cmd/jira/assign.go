package jira

import (
	"fmt"
	"strings"

	"github.com/joselrodrigues/atlassian/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var assignCmd = &cobra.Command{
	Use:   "assign [issue-key] [email-or-accountId]",
	Short: "Assign a user to an issue",
	Long: `Assign a user to an issue using email or accountId.
If the input contains '@', it's treated as an email and the accountId is looked up.
Use --unassign to remove the current assignee.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueKey := args[0]
		unassign, _ := cmd.Flags().GetBool("unassign")

		client := jira.NewClient()

		if unassign {
			if err := client.AssignIssue(issueKey, ""); err != nil {
				return fmt.Errorf("failed to unassign issue: %w", err)
			}

			baseURL := strings.TrimSuffix(viper.GetString("jira_base_url"), "/")
			fmt.Printf("Issue %s unassigned successfully!\n", issueKey)
			fmt.Printf("URL: %s/browse/%s\n", baseURL, issueKey)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("user email or accountId is required (or use --unassign)")
		}

		userInput := args[1]
		var accountID string

		if strings.Contains(userInput, "@") {
			users, err := client.SearchUsers(userInput)
			if err != nil {
				return fmt.Errorf("failed to search for user: %w", err)
			}

			if len(users) == 0 {
				return fmt.Errorf("no user found with email: %s", userInput)
			}

			accountID = users[0].AccountID
			fmt.Printf("Found user: %s (%s)\n", users[0].DisplayName, users[0].AccountID)
		} else {
			accountID = userInput
		}

		if err := client.AssignIssue(issueKey, accountID); err != nil {
			return fmt.Errorf("failed to assign issue: %w", err)
		}

		baseURL := strings.TrimSuffix(viper.GetString("jira_base_url"), "/")
		fmt.Printf("Issue %s assigned successfully!\n", issueKey)
		fmt.Printf("URL: %s/browse/%s\n", baseURL, issueKey)

		return nil
	},
}

func init() {
	Cmd.AddCommand(assignCmd)

	assignCmd.Flags().Bool("unassign", false, "Remove current assignee")
}
