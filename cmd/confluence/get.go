package confluence

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joselrodrigues/atlassian/internal/confluence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get [page-id]",
	Short: "Get a page by ID",
	Long:  `Retrieve a Confluence page by its ID, including metadata and optionally the body content.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pageID := args[0]
		bodyFormat, _ := cmd.Flags().GetString("body-format")
		output := viper.GetString("output")

		expand := []string{"version", "space"}
		if bodyFormat != "" {
			expand = append(expand, "body."+bodyFormat)
		}

		client := confluence.NewClient()
		page, err := client.GetPage(pageID, expand)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		printPage(page, output, bodyFormat)
	},
}

func init() {
	Cmd.AddCommand(getCmd)
	getCmd.Flags().String("body-format", "", "Include body content: storage, view")
}

func printPage(page *confluence.Page, format string, bodyFormat string) {
	if format == "json" {
		data, _ := json.MarshalIndent(page, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("| %-12s | %-60s |\n", "Field", "Value")
	fmt.Printf("| %-12s | %-60s |\n", "------------", "------------------------------------------------------------")
	fmt.Printf("| %-12s | %-60s |\n", "ID", page.ID)
	fmt.Printf("| %-12s | %-60s |\n", "Title", truncate(page.Title, 60))
	fmt.Printf("| %-12s | %-60s |\n", "Type", page.Type)
	fmt.Printf("| %-12s | %-60s |\n", "Status", page.Status)

	if page.Space != nil {
		fmt.Printf("| %-12s | %-60s |\n", "Space", page.Space.Key+" - "+page.Space.Name)
	}

	if page.Version != nil {
		fmt.Printf("| %-12s | %-60d |\n", "Version", page.Version.Number)
		if page.Version.By != nil {
			fmt.Printf("| %-12s | %-60s |\n", "Author", page.Version.By.DisplayName)
		}
		fmt.Printf("| %-12s | %-60s |\n", "Updated", page.Version.When)
	}

	fmt.Printf("| %-12s | %-60s |\n", "Web URL", page.Links.WebUI)

	if bodyFormat != "" && page.Body != nil {
		fmt.Println("\n--- Body Content ---")
		var content string
		if bodyFormat == "storage" && page.Body.Storage != nil {
			content = page.Body.Storage.Value
		} else if bodyFormat == "view" && page.Body.View != nil {
			content = page.Body.View.Value
		}
		if content != "" {
			fmt.Println(content)
		}
	}
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
