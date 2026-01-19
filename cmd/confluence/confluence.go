package confluence

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:     "confluence",
	Aliases: []string{"conf"},
	Short:   "Confluence operations",
	Long:    `Commands for interacting with Confluence: spaces, pages, search, and content management.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		validateConfig()
	},
}

func validateConfig() {
	if viper.GetString("confluence_token") == "" {
		fmt.Fprintln(os.Stderr, "Error: CONFLUENCE_TOKEN environment variable is required")
		os.Exit(1)
	}
	if viper.GetString("confluence_base_url") == "" {
		fmt.Fprintln(os.Stderr, "Error: CONFLUENCE_BASE_URL environment variable is required")
		os.Exit(1)
	}
}
