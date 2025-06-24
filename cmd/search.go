package cmd

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search for keys or projects",
	Long: `Search Secrets

Find keys or projects containing the specified keyword (case-insensitive).

EXAMPLES:
  uzp search api
  uzp search database
  uzp search myapp

OUTPUT:
  Shows matching projects and keys in the same format as list.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked, auto-unlock if needed
		if !vault.IsUnlocked() {
			fmt.Fprint(os.Stderr, "Enter master password: ")
			password, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Fprintln(os.Stderr) // New line after password

			if err := vault.Unlock(string(password)); err != nil {
				return fmt.Errorf("invalid password")
			}

			// Clear password from memory
			for i := range password {
				password[i] = 0
			}
		}

		keyword := args[0]

		// Search vault
		results, err := vault.Search(keyword)
		if err != nil {
			return err
		}

		// Check if no results
		if len(results) == 0 {
			fmt.Printf("No results found for '%s'\n", keyword)
			return nil
		}

		// Sort projects for consistent display
		projectNames := make([]string, 0, len(results))
		for project := range results {
			projectNames = append(projectNames, project)
		}
		sort.Strings(projectNames)

		// Display results
		fmt.Printf("Results for '%s':\n", keyword)

		for _, project := range projectNames {
			keys := results[project]
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
