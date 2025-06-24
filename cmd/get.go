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
	Long: `🔍 Retrieve Secret Value

Get a secret value from the vault by specifying the project and key.
The value will be printed to stdout (use copy command for clipboard).

FORMAT:
  project/key (e.g., myapp/api_key)

EXAMPLES:
  uzp get myapp/api_key              # Get API key for myapp
  uzp get backend/database_url       # Get database URL
  uzp get auth/jwt_secret            # Get JWT secret
  
💡 TIPS:
  • Use 'uzp list' to see all available secrets
  • Use 'uzp copy project/key' to copy to clipboard
  • Use 'uzp search keyword' to find specific secrets`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate arguments FIRST before prompting for password
		if len(args) == 0 {
			return fmt.Errorf(`❌ Missing secret path

USAGE:
  uzp get <project/key>

EXAMPLES:
  uzp get myapp/api_key
  uzp get backend/database_url

💡 TIP: Use 'uzp list' to see all available secrets`)
		}

		// Parse project/key
		parts := strings.Split(args[0], "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf(`❌ Invalid format: '%s'

REQUIRED FORMAT:
  project/key

EXAMPLES:
  uzp get myapp/api_key           ✅ Valid
  uzp get backend/database_url    ✅ Valid
  uzp get myapp                   ❌ Missing key
  uzp get /api_key                ❌ Missing project

💡 TIP: Use 'uzp list' to see all available secrets`, args[0])
		}

		// Check if vault is unlocked, auto-unlock if needed
		if !vault.IsUnlocked() {
			fmt.Fprint(os.Stderr, "🔒 Vault is locked. Enter master password: ")
			password, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Fprintln(os.Stderr) // New line after password

			if err := vault.Unlock(string(password)); err != nil {
				return fmt.Errorf("❌ Invalid master password. Please try again")
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
			return fmt.Errorf(`❌ Secret not found: %s/%s

💡 SUGGESTIONS:
  • Check spelling: uzp list (shows all secrets)
  • Search similar: uzp search %s
  • Add secret: uzp add

Available commands:
  uzp list                    # Show all secrets
  uzp search %s               # Search for '%s'
  uzp add                     # Add new secret`, project, key, key, key, key)
		}

		// Print value
		fmt.Println(value)

		return nil
	},
}
