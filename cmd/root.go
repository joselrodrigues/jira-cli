package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "jira-cli",
	Short: "CLI for interacting with Jira API",
	Long:  `A command-line interface for Jira operations including issues, comments, and transitions.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("output", "o", "text", "Output format: text, json")
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	viper.SetEnvPrefix("JIRA")
	viper.AutomaticEnv()

	if viper.GetString("TOKEN") == "" {
		fmt.Fprintln(os.Stderr, "Error: JIRA_TOKEN environment variable is required")
		os.Exit(1)
	}

	if viper.GetString("BASE_URL") == "" {
		fmt.Fprintln(os.Stderr, "Error: JIRA_BASE_URL environment variable is required")
		os.Exit(1)
	}
}
