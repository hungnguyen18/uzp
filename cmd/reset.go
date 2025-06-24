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
	Long:  `Reset the vault by deleting all stored secrets. This action is irreversible!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked
		if !vault.IsUnlocked() {
			return fmt.Errorf("vault is locked. Use 'uzp unlock' first")
		}

		// Confirm reset
		fmt.Println("⚠️  WARNING: This will permanently delete ALL secrets in the vault!")
		fmt.Println("This action CANNOT be undone.")
		fmt.Print("\nType 'DELETE ALL' to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		confirmation, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}
		confirmation = strings.TrimSpace(confirmation)

		if confirmation != "DELETE ALL" {
			fmt.Println("❌ Reset cancelled.")
			return nil
		}

		// Perform reset
		if err := vault.Reset(); err != nil {
			return fmt.Errorf("failed to reset vault: %w", err)
		}

		fmt.Println("✅ Vault has been reset. All data deleted.")
		fmt.Println("🔒 Vault is now locked.")

		return nil
	},
} 