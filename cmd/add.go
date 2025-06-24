package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a secret to the vault",
	Long: `➕ Add New Secret to Vault

Interactively add a new secret to the vault. You'll be prompted for:
  • Project name (groups related secrets)
  • Key name (identifier for the secret)
  • Value (the actual secret, hidden while typing)

ORGANIZATION:
  Secrets are organized by projects for better management.
  Example structure:
    myapp/
      ├── api_key
      ├── database_url
      └── jwt_secret

EXAMPLES:
  uzp add                           # Interactive mode
  
  # Example session:
  Project name: myapp
  Key name: api_key  
  Value (hidden): ****************
  ✅ Secret added successfully: myapp/api_key

💡 TIPS:
  • Use descriptive project names (e.g., 'backend', 'frontend', 'aws')
  • Use clear key names (e.g., 'api_key', 'database_url', 'jwt_secret')
  • Values are encrypted with AES-256-GCM before storage`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked, auto-unlock if needed
		if !vault.IsUnlocked() {
			fmt.Fprint(os.Stderr, "🔒 Vault is locked. Enter master password: ")
			password, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Fprintln(os.Stderr) // New line after password

			if err := vault.Unlock(string(password)); err != nil {
				return fmt.Errorf("❌ Invalid master password. Please try again")
			}

			// Clear password from memory
			for i := range password {
				password[i] = 0
			}

			fmt.Fprintln(os.Stderr, "✅ Vault unlocked successfully!")
		}

		reader := bufio.NewReader(os.Stdin)

		// Get project name
		fmt.Print("Project name: ")
		project, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read project name: %w", err)
		}
		project = strings.TrimSpace(project)

		// Get key name
		fmt.Print("Key name: ")
		key, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read key name: %w", err)
		}
		key = strings.TrimSpace(key)

		// Get value (sensitive, so use password input)
		fmt.Print("Value (hidden): ")
		valueBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read value: %w", err)
		}
		fmt.Println() // New line after password
		value := string(valueBytes)

		// Validate inputs
		if project == "" {
			return fmt.Errorf(`❌ Project name cannot be empty

💡 SUGGESTIONS:
  • Use descriptive names: 'myapp', 'backend', 'aws'
  • Group related secrets together
  • Use lowercase with hyphens: 'my-app', 'web-service'`)
		}

		if key == "" {
			return fmt.Errorf(`❌ Key name cannot be empty

💡 SUGGESTIONS:
  • Use descriptive names: 'api_key', 'database_url'
  • Use snake_case format: 'jwt_secret', 'oauth_token'
  • Be specific: 'stripe_api_key' vs 'api_key'`)
		}

		if value == "" {
			return fmt.Errorf(`❌ Value cannot be empty

💡 NOTE: The secret value must contain actual data to be stored`)
		}

		// Add to vault
		if err := vault.Add(project, key, value); err != nil {
			return fmt.Errorf("failed to add secret: %w", err)
		}

		fmt.Printf("✅ Secret added successfully: %s/%s\n", project, key)

		// Clear sensitive data from memory
		for i := range valueBytes {
			valueBytes[i] = 0
		}

		return nil
	},
}
