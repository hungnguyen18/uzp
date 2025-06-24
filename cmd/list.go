package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects and keys",
	Long:  `Display a list of all projects and their associated keys in the vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked
		if !vault.IsUnlocked() {
			return fmt.Errorf("vault is locked. Use 'uzp unlock' first")
		}

		// Get all projects and keys
		projects, err := vault.List()
		if err != nil {
			return err
		}

		// Check if vault is empty
		if len(projects) == 0 {
			fmt.Println("📭 Vault is empty. Use 'uzp add' to add secrets.")
			return nil
		}

		// Sort projects for consistent display
		projectNames := make([]string, 0, len(projects))
		for project := range projects {
			projectNames = append(projectNames, project)
		}
		sort.Strings(projectNames)

		// Display projects and keys
		fmt.Println("🔐 Vault Contents:")
		fmt.Println("==================")
		
		for _, project := range projectNames {
			keys := projects[project]
			sort.Strings(keys) // Sort keys for consistent display
			
			fmt.Printf("\n📁 %s\n", project)
			for _, key := range keys {
				fmt.Printf("   🔑 %s\n", key)
			}
		}

		return nil
	},
}