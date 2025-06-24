package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock the vault",
	Long:  `Unlock the vault by entering your master password. The vault will remain unlocked for the current session.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault exists
		if !vault.Exists() {
			return fmt.Errorf("vault does not exist. Use 'uzp init' to create one")
		}

		// Check if already unlocked
		if vault.IsUnlocked() {
			fmt.Println("🔓 Vault is already unlocked.")
			return nil
		}

		// Prompt for password
		fmt.Print("Enter master password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println() // New line after password

		// Unlock vault
		if err := vault.Unlock(string(password)); err != nil {
			return fmt.Errorf("failed to unlock vault: %w", err)
		}

		fmt.Println("✅ Vault unlocked successfully!")

		// Clear password from memory
		for i := range password {
			password[i] = 0
		}

		return nil
	},
} 