package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joselrodrigues/jira-cli/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long:  `Create a new Jira issue with the specified project, type, summary, and description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		issueType, _ := cmd.Flags().GetString("type")
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

		if summary == "" {
			return fmt.Errorf("--summary is required")
		}

		client := jira.NewClient()
		resp, err := client.CreateIssue(project, issueType, summary, description)
		if err != nil {
			return fmt.Errorf("failed to create issue: %w", err)
		}

		if viper.GetString("output") == "json" {
			data, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		baseURL := strings.TrimSuffix(viper.GetString("BASE_URL"), "/")
		fmt.Printf("Issue created successfully!\n")
		fmt.Printf("Key: %s\n", resp.Key)
		fmt.Printf("URL: %s/browse/%s\n", baseURL, resp.Key)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("project", "p", "", "Project key (required)")
	createCmd.MarkFlagRequired("project")
	createCmd.Flags().StringP("type", "t", "Story", "Issue type (Story, Bug, Task)")
	createCmd.Flags().StringP("summary", "s", "", "Issue summary (required)")
	createCmd.Flags().StringP("description", "d", "", "Issue description")
	createCmd.Flags().Bool("stdin", false, "Read description from stdin")
}
