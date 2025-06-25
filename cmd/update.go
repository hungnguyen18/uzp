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

var updateCmd = &cobra.Command{
	Use:   "update <project/key>",
	Short: "Update an existing secret value",
	Long: `Update Secret

Update the value of an existing secret in the vault.

FORMAT:
  project/key

EXAMPLES:
  uzp update myapp/api_key
  uzp update backend/database_url
  uzp update auth/jwt_secret

WORKFLOW:
  1. Specify secret path
  2. Enter new value (hidden)
  3. Confirm update`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate arguments FIRST before prompting for password
		if len(args) == 0 {
			return fmt.Errorf("usage: uzp update <project/key>")
		}

		// Parse project/key
		parts := strings.Split(args[0], "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid format. use: project/key")
		}

		project := parts[0]
		key := parts[1]

		// Check if vault is unlocked, prompt for password if needed
		if err := ensureVaultUnlocked(); err != nil {
			return err
		}

		// Check if secret exists
		currentValue, err := vault.Get(project, key)
		if err != nil {
			return fmt.Errorf("secret not found: %s/%s", project, key)
		}

		// Show current secret info (without value for security)
		fmt.Printf("Updating: %s/%s\n", project, key)

		// Get new value
		fmt.Print("New value: ")
		valueBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read value: %w", err)
		}
		fmt.Println() // New line after password
		newValue := string(valueBytes)

		// Validate new value
		if newValue == "" {
			return fmt.Errorf("value cannot be empty")
		}

		// Check if new value is the same as current
		if newValue == currentValue {
			fmt.Println("No changes made.")
			return nil
		}

		// Confirm update
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Update? (y/N): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}

		// Update secret
		if err := vault.Add(project, key, newValue); err != nil {
			return fmt.Errorf("failed to update secret: %w", err)
		}

		fmt.Printf("Updated: %s/%s\n", project, key)

		// Clear sensitive data from memory
		for i := range valueBytes {
			valueBytes[i] = 0
		}

		return nil
	},
}
