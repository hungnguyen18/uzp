package cmd

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

// ensureVaultUnlocked checks if vault is unlocked and prompts for password if needed
func ensureVaultUnlocked() error {
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
	return nil
}
