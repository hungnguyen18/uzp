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

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the vault (delete all data)",
	Long: `Reset Vault

Delete all stored secrets. This action is irreversible!

EXAMPLES:
  uzp reset

CONFIRMATION:
  You will be prompted to type 'DELETE ALL' to confirm.

WARNING:
  This permanently deletes all secrets and cannot be undone.`,
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

		// Confirm reset
		fmt.Println("WARNING: This will permanently delete ALL secrets in the vault!")
		fmt.Println("This action CANNOT be undone.")
		fmt.Print("\nType 'DELETE ALL' to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		confirmation, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}
		confirmation = strings.TrimSpace(confirmation)

		if confirmation != "DELETE ALL" {
			fmt.Println("Reset cancelled.")
			return nil
		}

		// Perform reset
		if err := vault.Reset(); err != nil {
			return fmt.Errorf("failed to reset vault: %w", err)
		}

		fmt.Println("Vault reset. All data deleted.")

		return nil
	},
}
