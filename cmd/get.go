package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <project/key>",
	Short: "Get a secret value from the vault",
	Long:  `Retrieve a secret value from the vault by specifying project/key.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if vault is unlocked
		if !vault.IsUnlocked() {
			return fmt.Errorf("vault is locked. Use 'uzp unlock' first")
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