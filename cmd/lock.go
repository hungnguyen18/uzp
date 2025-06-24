package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the vault",
	Long:  `Lock the vault manually. This will clear all decrypted data from memory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is already locked
		if !vault.IsUnlocked() {
			fmt.Println("🔒 Vault is already locked.")
			return nil
		}

		// Lock the vault
		vault.Lock()
		
		fmt.Println("✅ Vault locked successfully!")
		return nil
	},
} 