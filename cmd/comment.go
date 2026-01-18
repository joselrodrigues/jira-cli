package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/joselrodrigues/jira-cli/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage issue comments",
	Long:  `List or add comments to Jira issues.`,
}

var commentListCmd = &cobra.Command{
	Use:   "list [issue-key]",
	Short: "List comments on an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueKey := args[0]

		client := jira.NewClient()
		comments, err := client.GetComments(issueKey)
		if err != nil {
			return fmt.Errorf("failed to get comments: %w", err)
		}

		if viper.GetString("output") == "json" {
			data, _ := json.MarshalIndent(comments, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Comments on %s (%d total):\n\n", issueKey, comments.Total)
		for _, c := range comments.Comments {
			fmt.Printf("---\n")
			fmt.Printf("**%s** (%s)\n", c.Author.DisplayName, c.Created[:10])
			fmt.Printf("%s\n\n", c.Body)
		}

		return nil
	},
}

var commentAddCmd = &cobra.Command{
	Use:   "add [issue-key] [comment]",
	Short: "Add a comment to an issue",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueKey := args[0]
		body := args[1]

		client := jira.NewClient()
		comment, err := client.AddComment(issueKey, body)
		if err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}

		if viper.GetString("output") == "json" {
			data, _ := json.MarshalIndent(comment, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Comment added to %s (ID: %s)\n", issueKey, comment.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentAddCmd)
}