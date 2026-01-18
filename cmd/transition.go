package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/joserodrigues/jira-cli/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var transitionCmd = &cobra.Command{
	Use:   "transition",
	Short: "Manage issue transitions",
	Long:  `List available transitions or transition an issue to a new status.`,
}

var transitionListCmd = &cobra.Command{
	Use:   "list [issue-key]",
	Short: "List available transitions for an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueKey := args[0]

		client := jira.NewClient()
		transitions, err := client.GetTransitions(issueKey)
		if err != nil {
			return fmt.Errorf("failed to get transitions: %w", err)
		}

		if viper.GetString("output") == "json" {
			data, _ := json.MarshalIndent(transitions, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Available transitions for %s:\n\n", issueKey)
		fmt.Printf("| ID | Name | To Status |\n")
		fmt.Printf("|----|------|----------|\n")
		for _, t := range transitions.Transitions {
			fmt.Printf("| %s | %s | %s |\n", t.ID, t.Name, t.To.Name)
		}

		return nil
	},
}

var transitionDoCmd = &cobra.Command{
	Use:   "do [issue-key] [transition-name-or-id]",
	Short: "Transition an issue to a new status",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueKey := args[0]
		transition := args[1]

		client := jira.NewClient()
		if err := client.DoTransition(issueKey, transition); err != nil {
			return fmt.Errorf("failed to transition: %w", err)
		}

		fmt.Printf("Issue %s transitioned to '%s' successfully!\n", issueKey, transition)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(transitionCmd)
	transitionCmd.AddCommand(transitionListCmd)
	transitionCmd.AddCommand(transitionDoCmd)
}
