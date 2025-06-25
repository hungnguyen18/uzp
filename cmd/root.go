package cmd

import (
	"fmt"

	"github.com/hungnguyen18/uzp-cli/internal/storage"
	"github.com/spf13/cobra"
)

// Version information - injected at build time
var Version = "dev" // Default value, will be overridden by ldflags during build

var (
	vault       *storage.Vault
	showVersion bool
	rootCmd     = &cobra.Command{
		Use:   "uzp",
		Short: "Secure secrets manager",
		Long: `UZP - Secure Secrets Manager

A command-line tool for securely storing and managing sensitive information
such as API keys, access tokens, and service credentials.

SECURITY:
  AES-256-GCM encryption, master password protection, local storage

BASIC USAGE:
  uzp init                    Initialize new vault
  uzp add                     Add secret
  uzp list                    List all secrets
  uzp get project/key         Get secret value
  uzp update project/key      Update secret
  uzp inject -p project       Export as environment variables

EXAMPLES:
  uzp inject -p myapp > .env  Export secrets to .env file
  uzp copy myapp/api_key      Copy secret to clipboard
  uzp search database         Search for secrets

STORAGE: ~/.uzp/uzp.vault (encrypted)`,
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("uzp version %s\n", Version)
				return
			}
			// If no subcommand is provided, show help by default
			cmd.Help()
		},
	}
)

func init() {
	// Add version flags
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version information")

	// Initialize vault instance
	vault = storage.NewVault()

	// Add all subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
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
