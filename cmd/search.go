package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search for keys or projects",
	Long:  `Search for keys or projects containing the specified keyword (case-insensitive).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked
		if !vault.IsUnlocked() {
			return fmt.Errorf("vault is locked. Use 'uzp unlock' first")
		}

		keyword := args[0]

		// Search vault
		results, err := vault.Search(keyword)
		if err != nil {
			return err
		}

		// Check if no results
		if len(results) == 0 {
			fmt.Printf("❌ No results found for '%s'\n", keyword)
			return nil
		}

		// Sort projects for consistent display
		projectNames := make([]string, 0, len(results))
		for project := range results {
			projectNames = append(projectNames, project)
		}
		sort.Strings(projectNames)

		// Display results
		fmt.Printf("🔍 Search results for '%s':\n", keyword)
		fmt.Println("==========================")
		
		for _, project := range projectNames {
			keys := results[project]
			sort.Strings(keys) // Sort keys for consistent display
			
			fmt.Printf("\n📁 %s\n", project)
			for _, key := range keys {
				fmt.Printf("   🔑 %s\n", key)
			}
		}

		return nil
	},
} 