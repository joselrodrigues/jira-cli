package cmd

import (
	"fmt"
	"os"

	"github.com/joselrodrigues/atlassian/cmd/confluence"
	"github.com/joselrodrigues/atlassian/cmd/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "atlassian",
	Short: "CLI for interacting with Atlassian products (Jira, Confluence)",
	Long:  `A command-line interface for Atlassian products including Jira operations (issues, comments, transitions) and Confluence (spaces, pages, search).`,
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

	rootCmd.AddCommand(jira.Cmd)
	rootCmd.AddCommand(confluence.Cmd)
}

func initConfig() {
	viper.AutomaticEnv()

	viper.BindEnv("jira_token", "JIRA_TOKEN")
	viper.BindEnv("jira_base_url", "JIRA_BASE_URL")
	viper.BindEnv("confluence_token", "CONFLUENCE_TOKEN")
	viper.BindEnv("confluence_base_url", "CONFLUENCE_BASE_URL")
}
