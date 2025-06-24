package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/hungnguyen/uzp/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	ttl int
)

var copyCmd = &cobra.Command{
	Use:   "copy <project/key>",
	Short: "Copy a secret value to clipboard",
	Long:  `Copy a secret value to clipboard. The value will be automatically cleared after TTL (default 15 seconds).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked, auto-unlock if needed
		if !vault.IsUnlocked() {
			fmt.Fprint(os.Stderr, "Vault is locked. Enter master password: ")
			password, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Fprintln(os.Stderr) // New line after password

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

		// Copy to clipboard
		duration := time.Duration(ttl) * time.Second
		if err := utils.CopyToClipboard(value, duration); err != nil {
			return err
		}

		fmt.Printf("✅ Copied %s/%s to clipboard.\n", project, key)
		fmt.Printf("📋 Clipboard will be cleared in %d seconds.\n", ttl)

		return nil
	},
}

func init() {
	copyCmd.Flags().IntVarP(&ttl, "ttl", "t", 15, "Time to live in seconds before clipboard is cleared")
}
