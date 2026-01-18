package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joserodrigues/jira-cli/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update [issue-key]",
	Short: "Update an existing issue",
	Long:  `Update fields of an existing Jira issue.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueKey := args[0]
		summary, _ := cmd.Flags().GetString("summary")
		description, _ := cmd.Flags().GetString("description")
		fromStdin, _ := cmd.Flags().GetBool("stdin")

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

		fields := make(map[string]interface{})
		if summary != "" {
			fields["summary"] = summary
		}
		if description != "" {
			fields["description"] = description
		}

		if len(fields) == 0 {
			return fmt.Errorf("at least one field must be specified (--summary or --description)")
		}

		client := jira.NewClient()
		if err := client.UpdateIssue(issueKey, fields); err != nil {
			return fmt.Errorf("failed to update issue: %w", err)
		}

		baseURL := strings.TrimSuffix(viper.GetString("BASE_URL"), "/")
		fmt.Printf("Issue %s updated successfully!\n", issueKey)
		fmt.Printf("URL: %s/browse/%s\n", baseURL, issueKey)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("summary", "s", "", "New summary")
	updateCmd.Flags().StringP("description", "d", "", "New description")
	updateCmd.Flags().Bool("stdin", false, "Read description from stdin")
}