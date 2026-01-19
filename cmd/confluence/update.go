package confluence

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joselrodrigues/atlassian/internal/confluence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update [page-id]",
	Short: "Update an existing page",
	Long: `Update a Confluence page's title or content.

Examples:
  atlassian confluence update 123456 --title "New Title"
  echo "<p>New content</p>" | atlassian confluence update 123456 --stdin
  atlassian confluence update 123456 --title "Title" --message "Updated via CLI"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pageID := args[0]
		title, _ := cmd.Flags().GetString("title")
		message, _ := cmd.Flags().GetString("message")
		useStdin, _ := cmd.Flags().GetBool("stdin")
		output := viper.GetString("output")

		client := confluence.NewClient()

		currentPage, err := client.GetPage(pageID, []string{"body.storage", "version"})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current page: %v\n", err)
			os.Exit(1)
		}

		newTitle := currentPage.Title
		if title != "" {
			newTitle = title
		}

		var newBody string
		if useStdin {
			newBody = readStdin()
		} else if currentPage.Body != nil && currentPage.Body.Storage != nil {
			newBody = currentPage.Body.Storage.Value
		}

		if title == "" && !useStdin {
			fmt.Fprintln(os.Stderr, "Error: must specify --title or --stdin")
			os.Exit(1)
		}

		currentVersion := 1
		if currentPage.Version != nil {
			currentVersion = currentPage.Version.Number
		}

		page, err := client.UpdatePage(pageID, newTitle, newBody, currentVersion, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(page, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Printf("Page updated successfully!\n")
			fmt.Printf("ID: %s\n", page.ID)
			fmt.Printf("Title: %s\n", page.Title)
			fmt.Printf("Version: %d\n", page.Version.Number)
			fmt.Printf("URL: %s%s\n", viper.GetString("confluence_base_url"), page.Links.WebUI)
		}
	},
}

func init() {
	Cmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("title", "t", "", "New page title")
	updateCmd.Flags().String("message", "", "Version message")
	updateCmd.Flags().Bool("stdin", false, "Read new body content from stdin")
}
