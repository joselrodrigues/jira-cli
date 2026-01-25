package jira

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/joselrodrigues/atlassian/internal/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldsCmd = &cobra.Command{
	Use:   "fields",
	Short: "List Jira fields",
	Long:  `List all Jira fields to discover custom field IDs (e.g., Story Points, Sprint).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		customOnly, _ := cmd.Flags().GetBool("custom")
		output := viper.GetString("output")

		client := jira.NewClient()
		fields, err := client.GetFields()
		if err != nil {
			return fmt.Errorf("failed to get fields: %w", err)
		}

		var filtered []jira.Field
		for _, field := range fields {
			if customOnly && !field.Custom {
				continue
			}
			if name != "" && !strings.Contains(strings.ToLower(field.Name), strings.ToLower(name)) {
				continue
			}
			filtered = append(filtered, field)
		}

		if output == "json" {
			jsonData, err := json.MarshalIndent(filtered, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(jsonData))
			return nil
		}

		if len(filtered) == 0 {
			fmt.Println("No fields found matching criteria")
			return nil
		}

		fmt.Println("| Field ID | Name | Custom | Type |")
		fmt.Println("| -------- | ---- | ------ | ---- |")
		for _, field := range filtered {
			custom := "No"
			if field.Custom {
				custom = "Yes"
			}
			fieldType := field.Schema.Type
			if fieldType == "" {
				fieldType = "-"
			}
			fmt.Printf("| %s | %s | %s | %s |\n", field.ID, field.Name, custom, fieldType)
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(fieldsCmd)

	fieldsCmd.Flags().StringP("name", "n", "", "Filter fields by name (case-insensitive)")
	fieldsCmd.Flags().Bool("custom", false, "Show only custom fields")
}
