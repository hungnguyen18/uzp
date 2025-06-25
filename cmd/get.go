package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
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

		// Check if vault is unlocked, prompt for password if needed
		if err := ensureVaultUnlocked(); err != nil {
			return err
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
