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
	Long:  `Add a new secret to the vault by specifying project, key, and value.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked, auto-unlock if needed
		if !vault.IsUnlocked() {
			fmt.Print("Vault is locked. Enter master password: ")
			password, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Println() // New line after password

			if err := vault.Unlock(string(password)); err != nil {
				return fmt.Errorf("failed to unlock vault: %w", err)
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
		if project == "" || key == "" || value == "" {
			return fmt.Errorf("project, key, and value cannot be empty")
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
