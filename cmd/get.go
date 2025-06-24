package cmd

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var getCmd = &cobra.Command{
	Use:   "get <project/key>",
	Short: "Get a secret value from the vault",
	Long:  `Retrieve a secret value from the vault by specifying project/key.`,
	Args:  cobra.ExactArgs(1),
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

		// Parse project/key
		parts := strings.Split(args[0], "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format. Use: project/key")
		}

		project := parts[0]
		key := parts[1]

		// Get value
		value, err := vault.Get(project, key)
		if err != nil {
			return err
		}

		// Print value
		fmt.Println(value)

		return nil
	},
}
