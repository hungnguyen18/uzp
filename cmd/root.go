package cmd

import (
	"github.com/hungnguyen/uzp/internal/storage"
	"github.com/spf13/cobra"
)

var (
	vault   *storage.Vault
	rootCmd = &cobra.Command{
		Use:   "uzp",
		Short: "A secure CLI tool for managing secrets",
		Long: `🔐 UZP - Secure Secrets Manager

uzp is a command-line tool designed to securely store and manage sensitive 
information such as API keys, access tokens, and service credentials.

SECURITY FEATURES:
  • AES-256-GCM encryption for maximum security
  • Master password protection with scrypt key derivation  
  • Local storage with 0600 file permissions
  • Automatic clipboard clearing after TTL

GETTING STARTED:
  1. Initialize a new vault:    uzp init
  2. Add your first secret:    uzp add
  3. List all secrets:         uzp list
  4. Get a secret value:       uzp get project/key
  5. Export to .env file:      uzp inject -p project > .env

EXAMPLES:
  uzp init                           # Create new vault
  uzp add                            # Add secret interactively  
  uzp get myapp/api_key             # Get secret value
  uzp copy myapp/api_key            # Copy to clipboard
  uzp list                          # Show all secrets
  uzp search api                    # Search for secrets
  uzp inject -p myapp > .env        # Export as .env file
  uzp lock                          # Lock the vault
  uzp unlock                        # Unlock the vault

STORAGE LOCATION:
  ~/.uzp/uzp.vault (encrypted with your master password)

For detailed help on any command, use: uzp [command] --help`,
	}
)

func init() {
	// Initialize vault instance
	vault = storage.NewVault()

	// Add all subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(unlockCmd)
	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(injectCmd)
	rootCmd.AddCommand(resetCmd)
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
