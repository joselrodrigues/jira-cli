package confluence

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joselrodrigues/atlassian/internal/confluence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new page",
	Long: `Create a new Confluence page in a space.

Examples:
  atlassian confluence create -s MYSPACE -t "My New Page"
  atlassian confluence create -s MYSPACE -t "Child Page" --parent 123456
  echo "<p>Content</p>" | atlassian confluence create -s MYSPACE -t "Page" --stdin`,
	Run: func(cmd *cobra.Command, args []string) {
		spaceKey, _ := cmd.Flags().GetString("space")
		title, _ := cmd.Flags().GetString("title")
		parentID, _ := cmd.Flags().GetString("parent")
		useStdin, _ := cmd.Flags().GetBool("stdin")
		output := viper.GetString("output")

		if spaceKey == "" {
			fmt.Fprintln(os.Stderr, "Error: --space/-s flag is required")
			os.Exit(1)
		}
		if title == "" {
			fmt.Fprintln(os.Stderr, "Error: --title/-t flag is required")
			os.Exit(1)
		}

		var body string
		if useStdin {
			body = readStdin()
		} else {
			body = "<p>New page content</p>"
		}

		client := confluence.NewClient()
		page, err := client.CreatePage(spaceKey, title, body, parentID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(page, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Printf("Page created successfully!\n")
			fmt.Printf("ID: %s\n", page.ID)
			fmt.Printf("Title: %s\n", page.Title)
			fmt.Printf("URL: %s%s\n", viper.GetString("confluence_base_url"), page.Links.WebUI)
		}
	},
}

func init() {
	Cmd.AddCommand(createCmd)
	createCmd.Flags().StringP("space", "s", "", "Space key (required)")
	createCmd.Flags().StringP("title", "t", "", "Page title (required)")
	createCmd.Flags().String("parent", "", "Parent page ID")
	createCmd.Flags().Bool("stdin", false, "Read body content from stdin")
}

func readStdin() string {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}
