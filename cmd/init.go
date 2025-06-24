package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault",
	Long:  `Initialize a new vault with a master password. This command creates an encrypted vault file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault already exists
		if vault.Exists() {
			return fmt.Errorf("vault already exists. Use 'uzp unlock' to access it")
		}

		// Prompt for master password
		fmt.Print("Enter master password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println() // New line after password

		// Confirm password
		fmt.Print("Confirm master password: ")
		confirmPassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println() // New line after password

		// Check if passwords match
		if string(password) != string(confirmPassword) {
			return fmt.Errorf("passwords do not match")
		}

		// Check password strength
		if len(password) < 8 {
			return fmt.Errorf("password must be at least 8 characters long")
		}

		// Initialize vault
		if err := vault.Initialize(string(password)); err != nil {
			return fmt.Errorf("failed to initialize vault: %w", err)
		}

		fmt.Println("✅ Vault initialized successfully!")
		fmt.Println("🔓 Vault is now unlocked and ready to use.")
		
		// Clear password from memory
		for i := range password {
			password[i] = 0
		}
		for i := range confirmPassword {
			confirmPassword[i] = 0
		}

		return nil
	},
} 