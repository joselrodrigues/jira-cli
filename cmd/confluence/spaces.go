package confluence

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joselrodrigues/atlassian/internal/confluence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var spacesCmd = &cobra.Command{
	Use:   "spaces [space-key]",
	Short: "List spaces or get space details",
	Long:  `List all available Confluence spaces or get details of a specific space by its key.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := confluence.NewClient()
		output := viper.GetString("output")

		if len(args) == 1 {
			space, err := client.GetSpace(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			printSpace(space, output)
		} else {
			limit, _ := cmd.Flags().GetInt("limit")
			spaces, err := client.ListSpaces(limit)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			printSpaces(spaces, output)
		}
	},
}

func init() {
	Cmd.AddCommand(spacesCmd)
	spacesCmd.Flags().Int("limit", 25, "Maximum number of spaces to return")
}

func printSpace(space *confluence.Space, format string) {
	if format == "json" {
		data, _ := json.MarshalIndent(space, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("| %-12s | %-50s |\n", "Field", "Value")
	fmt.Printf("| %-12s | %-50s |\n", "------------", "--------------------------------------------------")
	fmt.Printf("| %-12s | %-50d |\n", "ID", space.ID)
	fmt.Printf("| %-12s | %-50s |\n", "Key", space.Key)
	fmt.Printf("| %-12s | %-50s |\n", "Name", space.Name)
	fmt.Printf("| %-12s | %-50s |\n", "Status", space.Status)
	fmt.Printf("| %-12s | %-50s |\n", "Type", space.Type)
}

func printSpaces(spaces *confluence.SpacesResponse, format string) {
	if format == "json" {
		data, _ := json.MarshalIndent(spaces, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("| %-10s | %-40s | %-10s |\n", "Key", "Name", "Type")
	fmt.Printf("| %-10s | %-40s | %-10s |\n", "----------", "----------------------------------------", "----------")
	for _, s := range spaces.Results {
		name := s.Name
		if len(name) > 40 {
			name = name[:37] + "..."
		}
		fmt.Printf("| %-10s | %-40s | %-10s |\n", s.Key, name, s.Type)
	}
	fmt.Printf("\nTotal: %d spaces\n", spaces.Size)
}
