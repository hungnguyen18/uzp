package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects and keys",
	Long: `List Secrets

Display all projects and their keys in the vault.

EXAMPLES:
  uzp list

OUTPUT FORMAT:
  project1:
    key1
    key2
  
  project2:
    key3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked, prompt for password if needed
		if err := ensureVaultUnlocked(); err != nil {
			return err
		}

		// Get all projects and keys
		projects, err := vault.List()
		if err != nil {
			return err
		}

		// Check if vault is empty
		if len(projects) == 0 {
			fmt.Println("No secrets found.")
			return nil
		}

		// Sort projects for consistent display
		projectNames := make([]string, 0, len(projects))
		for project := range projects {
			projectNames = append(projectNames, project)
		}
		sort.Strings(projectNames)

		// Display projects and keys
		for _, project := range projectNames {
			keys := projects[project]
			sort.Strings(keys) // Sort keys for consistent display

			fmt.Printf("%s:\n", project)
			for _, key := range keys {
				fmt.Printf("  %s\n", key)
			}
			fmt.Println()
		}

		return nil
	},
}
