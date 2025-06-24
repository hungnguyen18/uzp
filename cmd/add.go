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
	Long: `Add Secret

Interactively add a new secret or update an existing one.

WORKFLOW:
  1. Enter project name (groups related secrets)
  2. Enter key name (identifier for the secret)  
  3. Enter value (hidden while typing)
  4. If key exists, confirm update

ORGANIZATION:
  Secrets are organized by projects:
    myapp/api_key
    myapp/database_url
    backend/jwt_secret

EXAMPLES:
  uzp add                 Interactive mode

  Project name: myapp
  Key name: api_key
  Value (hidden): 
  Added: myapp/api_key`,
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
			return fmt.Errorf("Project name cannot be empty")
		}

		if key == "" {
			return fmt.Errorf("Key name cannot be empty")
		}

		if value == "" {
			return fmt.Errorf("Value cannot be empty")
		}

		// Check if key already exists
		existingValue, err := vault.Get(project, key)
		isUpdate := err == nil && existingValue != ""

		if isUpdate {
			fmt.Printf("Secret '%s/%s' already exists.\n", project, key)
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
		}

		// Add/Update to vault
		if err := vault.Add(project, key, value); err != nil {
			if isUpdate {
				return fmt.Errorf("failed to update secret: %w", err)
			}
			return fmt.Errorf("failed to add secret: %w", err)
		}

		if isUpdate {
			fmt.Printf("Updated: %s/%s\n", project, key)
		} else {
			fmt.Printf("Added: %s/%s\n", project, key)
		}

		// Clear sensitive data from memory
		for i := range valueBytes {
			valueBytes[i] = 0
		}

		return nil
	},
}
