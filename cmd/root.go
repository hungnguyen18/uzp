package cmd

import (
	"github.com/spf13/cobra"
	"github.com/hungnguyen/uzp/internal/storage"
)

var (
	vault *storage.Vault
	rootCmd = &cobra.Command{
		Use:   "uzp",
		Short: "A CLI tool for managing secrets securely",
		Long: `uzp is a command-line tool designed to securely store and manage 
sensitive information such as API keys, access tokens, and service credentials.

All data is encrypted using AES-256-GCM and stored locally.`,
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