package jira

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:   "jira",
	Short: "Jira operations",
	Long:  `Commands for interacting with Jira: issues, comments, transitions, sprints, and more.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		validateConfig()
	},
}

func validateConfig() {
	if viper.GetString("jira_token") == "" {
		fmt.Fprintln(os.Stderr, "Error: JIRA_TOKEN environment variable is required")
		os.Exit(1)
	}
	if viper.GetString("jira_base_url") == "" {
		fmt.Fprintln(os.Stderr, "Error: JIRA_BASE_URL environment variable is required")
		os.Exit(1)
	}
}
