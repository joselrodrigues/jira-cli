package jira

import (
	"encoding/json"
	"fmt"

	"github.com/joselrodrigues/atlassian/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Search for users",
	Long:  `Search for users by name or email to get their accountId.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		output := viper.GetString("output")

		if query == "" {
			return fmt.Errorf("--query is required")
		}

		client := jira.NewClient()
		users, err := client.SearchUsers(query)
		if err != nil {
			return fmt.Errorf("failed to search users: %w", err)
		}

		if output == "json" {
			jsonData, err := json.MarshalIndent(users, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(jsonData))
			return nil
		}

		if len(users) == 0 {
			fmt.Println("No users found")
			return nil
		}

		fmt.Println("| Account ID | Display Name | Email | Active |")
		fmt.Println("| ---------- | ------------ | ----- | ------ |")
		for _, user := range users {
			email := user.EmailAddress
			if email == "" {
				email = "-"
			}
			active := "Yes"
			if !user.Active {
				active = "No"
			}
			fmt.Printf("| %s | %s | %s | %s |\n", user.AccountID, user.DisplayName, email, active)
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(usersCmd)

	usersCmd.Flags().StringP("query", "q", "", "Search query (name or email)")
}
