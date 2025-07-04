package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var projectName string

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Export project secrets as .env format",
	Long: `Export Secrets

Export all secrets for a project in environment variable format.

USAGE:
  uzp inject --project PROJECT_NAME [> output_file]

EXAMPLES:
  uzp inject -p myapp > .env          Export to .env file
  uzp inject -p myapp > .env.example  Export to example file
  uzp inject -p backend               Display to terminal

WORKFLOW:
  1. Specify project with -p flag
  2. Optionally redirect output to file
  3. Environment variables are formatted automatically

OUTPUT FORMAT:
  # Environment variables for project: myapp
  API_KEY=your_secret_value
  DATABASE_URL=your_connection_string

NOTE:
  Keys are converted to UPPERCASE with underscores`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate arguments FIRST - show help immediately if missing project
		if projectName == "" {
			return fmt.Errorf("missing project name\n\nusage: uzp inject -p PROJECT_NAME > .env\n\nSee 'uzp inject --help' for examples")
		}

		// Check if vault is unlocked, prompt for password if needed
		if err := ensureVaultUnlocked(); err != nil {
			return err
		}

		// Get project secrets
		secrets, err := vault.GetProjectSecrets(projectName)
		if err != nil {
			return fmt.Errorf("project not found: %s", projectName)
		}

		// Sort keys for consistent output
		keys := make([]string, 0, len(secrets))
		for key := range secrets {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// Show success feedback to stderr (won't be redirected to file)
		fmt.Fprintf(os.Stderr, "Exporting %d secrets from project '%s'\n", len(secrets), projectName)

		// Output in .env format
		fmt.Printf("# Environment variables for project: %s\n", projectName)
		fmt.Printf("# Generated by uzp\n\n")

		for _, key := range keys {
			// Convert key to uppercase and replace non-alphanumeric chars with underscore
			envKey := convertToEnvKey(key)
			fmt.Printf("%s=%s\n", envKey, secrets[key])
		}

		// Success message to stderr
		fmt.Fprintf(os.Stderr, "Successfully exported environment variables\n")

		return nil
	},
}

func init() {
	injectCmd.Flags().StringVarP(&projectName, "project", "p", "", "Project name to export secrets from")
}

// convertToEnvKey converts a key to environment variable format
func convertToEnvKey(key string) string {
	result := make([]byte, 0, len(key))

	for i := 0; i < len(key); i++ {
		c := key[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			// Convert to uppercase if lowercase
			if c >= 'a' && c <= 'z' {
				c = c - 32
			}
			result = append(result, c)
		} else {
			// Replace non-alphanumeric with underscore
			result = append(result, '_')
		}
	}

	return string(result)
}
