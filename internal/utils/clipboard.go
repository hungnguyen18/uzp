package utils

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

// CopyToClipboard copies text to clipboard and clears it after TTL
func CopyToClipboard(text string, ttl time.Duration) error {
	// Copy to clipboard
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	// Clear clipboard after TTL
	go func() {
		time.Sleep(ttl)
		// Clear by writing empty string
		clipboard.WriteAll("")
	}()

	return nil
} 