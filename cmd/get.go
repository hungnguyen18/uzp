package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var getCmd = &cobra.Command{
	Use:   "get <project/key>",
	Short: "Get a secret value from the vault",
	Long: `Get Secret

Retrieve a secret value from the vault and print to stdout.

FORMAT:
  project/key

EXAMPLES:
  uzp get myapp/api_key
  uzp get backend/database_url
  uzp get auth/jwt_secret`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate arguments FIRST before prompting for password
		if len(args) == 0 {
			return fmt.Errorf("usage: uzp get <project/key>")
		}

		// Parse project/key
		parts := strings.Split(args[0], "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid format. use: project/key")
		}

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

		project := parts[0]
		key := parts[1]

		// Get value
		value, err := vault.Get(project, key)
		if err != nil {
			return fmt.Errorf("secret not found: %s/%s", project, key)
		}

		// Print value
		fmt.Println(value)

		return nil
	},
}
