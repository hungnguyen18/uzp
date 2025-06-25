package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
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
		// Check if vault is unlocked, prompt for password if needed
		if err := ensureVaultUnlocked(); err != nil {
			return err
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
