package storage

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hungnguyen18/uzp-cli/internal/crypto"
)

type VaultData struct {
	Version  int                          `json:"version"`
	Salt     string                       `json:"salt"`
	Hash     string                       `json:"hash"` // Password hash for verification
	Projects map[string]map[string]string `json:"projects"`
}

type EncryptedVault struct {
	Salt string `json:"salt"`
	Hash string `json:"hash"`
	Data string `json:"data"` // Base64 encoded encrypted data
}

// Session data for persistence
type SessionData struct {
	UnlockedAt time.Time `json:"unlocked_at"`
	KeyHash    string    `json:"key_hash"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type Vault struct {
	path        string
	sessionPath string
	data        *VaultData
	key         []byte
	unlocked    bool
	unlockedAt  time.Time
}

// Session timeout - 15 minutes
const SessionTimeout = 15 * time.Minute

// NewVault creates a new vault instance
func NewVault() *Vault {
	homeDir, _ := os.UserHomeDir()
	vaultDir := filepath.Join(homeDir, ".uzp")
	vaultPath := filepath.Join(vaultDir, "uzp.vault")
	sessionPath := filepath.Join(vaultDir, ".session")

	v := &Vault{
		path:        vaultPath,
		sessionPath: sessionPath,
	}

	// Try to restore session on creation
	v.tryRestoreSession()

	return v
}

// tryRestoreSession attempts to restore session from file
func (v *Vault) tryRestoreSession() {
	sessionData, err := v.loadSession()
	if err != nil {
		return // No valid session found
	}

	// Check if session is expired
	if time.Now().After(sessionData.ExpiresAt) {
		v.clearSession()
		return
	}

	// Try to unlock vault using saved session
	if err := v.restoreFromSession(sessionData); err != nil {
		v.clearSession()
		return
	}
}

// saveSession saves current session to file
func (v *Vault) saveSession() error {
	if !v.unlocked {
		return nil
	}

	sessionData := SessionData{
		UnlockedAt: v.unlockedAt,
		KeyHash:    crypto.HashData(v.key),
		ExpiresAt:  time.Now().Add(SessionTimeout),
	}

	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	// Create vault directory if not exists
	dir := filepath.Dir(v.sessionPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(v.sessionPath, data, 0600)
}

// loadSession loads session from file
func (v *Vault) loadSession() (*SessionData, error) {
	data, err := os.ReadFile(v.sessionPath)
	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	if err := json.Unmarshal(data, &sessionData); err != nil {
		return nil, err
	}

	return &sessionData, nil
}

// restoreFromSession restores vault state from session data
func (v *Vault) restoreFromSession(sessionData *SessionData) error {
	// We need to derive the key again - we can't safely store it in session
	// For now, we'll just mark as needing unlock
	v.unlocked = false
	return nil
}

// clearSession removes session file
func (v *Vault) clearSession() {
	os.Remove(v.sessionPath)
}

// Initialize creates a new vault with the given master password
func (v *Vault) Initialize(masterPassword string) error {
	// Create vault directory if not exists
	dir := filepath.Dir(v.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	// Check if vault already exists
	if _, err := os.Stat(v.path); err == nil {
		return fmt.Errorf("vault already exists at %s", v.path)
	}

	// Generate salt
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return err
	}

	// Derive key from password
	key, err := crypto.DeriveKey(masterPassword, salt)
	if err != nil {
		return err
	}

	// Create initial vault data
	v.data = &VaultData{
		Version:  1,
		Salt:     base64.StdEncoding.EncodeToString(salt),
		Hash:     crypto.HashPassword(masterPassword),
		Projects: make(map[string]map[string]string),
	}

	v.key = key
	v.unlocked = true
	v.unlockedAt = time.Now()

	// Save vault and session
	if err := v.save(); err != nil {
		return err
	}
	v.saveSession()

	return nil
}

// Unlock unlocks the vault with the master password
func (v *Vault) Unlock(masterPassword string) error {
	// Load encrypted vault
	encVault, err := v.loadEncrypted()
	if err != nil {
		return err
	}

	// Verify password hash
	if crypto.HashPassword(masterPassword) != encVault.Hash {
		return fmt.Errorf("invalid master password")
	}

	// Decode salt
	salt, err := base64.StdEncoding.DecodeString(encVault.Salt)
	if err != nil {
		return fmt.Errorf("failed to decode salt: %w", err)
	}

	// Derive key
	key, err := crypto.DeriveKey(masterPassword, salt)
	if err != nil {
		return err
	}

	// Decrypt data
	encryptedData, err := base64.StdEncoding.DecodeString(encVault.Data)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	decryptedData, err := crypto.Decrypt(encryptedData, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt vault: %w", err)
	}

	// Unmarshal vault data
	var vaultData VaultData
	if err := json.Unmarshal(decryptedData, &vaultData); err != nil {
		return fmt.Errorf("failed to unmarshal vault data: %w", err)
	}

	v.data = &vaultData
	v.key = key
	v.unlocked = true
	v.unlockedAt = time.Now()

	// Save session
	v.saveSession()

	return nil
}

// Lock locks the vault
func (v *Vault) Lock() {
	v.key = nil
	v.unlocked = false
	v.data = nil
	v.clearSession()
}

// IsUnlocked checks if vault is unlocked
func (v *Vault) IsUnlocked() bool {
	// Check session validity if not currently unlocked
	if !v.unlocked {
		v.tryRestoreSession()
	}
	return v.unlocked
}

// Add adds a secret to the vault
func (v *Vault) Add(project, key, value string) error {
	if !v.unlocked {
		return fmt.Errorf("vault is locked")
	}

	if v.data.Projects[project] == nil {
		v.data.Projects[project] = make(map[string]string)
	}

	v.data.Projects[project][key] = value
	return v.save()
}

// Get retrieves a secret from the vault
func (v *Vault) Get(project, key string) (string, error) {
	if !v.unlocked {
		return "", fmt.Errorf("vault is locked")
	}

	if proj, ok := v.data.Projects[project]; ok {
		if val, ok := proj[key]; ok {
			return val, nil
		}
	}

	return "", fmt.Errorf("secret not found: %s/%s", project, key)
}

// List returns all projects and keys
func (v *Vault) List() (map[string][]string, error) {
	if !v.unlocked {
		return nil, fmt.Errorf("vault is locked")
	}

	result := make(map[string][]string)
	for project, secrets := range v.data.Projects {
		keys := make([]string, 0, len(secrets))
		for key := range secrets {
			keys = append(keys, key)
		}
		result[project] = keys
	}

	return result, nil
}

// Search searches for keys or projects containing the keyword
func (v *Vault) Search(keyword string) (map[string][]string, error) {
	if !v.unlocked {
		return nil, fmt.Errorf("vault is locked")
	}

	result := make(map[string][]string)
	for project, secrets := range v.data.Projects {
		matches := []string{}
		for key := range secrets {
			if contains(project, keyword) || contains(key, keyword) {
				matches = append(matches, key)
			}
		}
		if len(matches) > 0 {
			result[project] = matches
		}
	}

	return result, nil
}

// GetProjectSecrets returns all secrets for a project
func (v *Vault) GetProjectSecrets(project string) (map[string]string, error) {
	if !v.unlocked {
		return nil, fmt.Errorf("vault is locked")
	}

	if proj, ok := v.data.Projects[project]; ok {
		// Return a copy to prevent external modification
		result := make(map[string]string)
		for k, v := range proj {
			result[k] = v
		}
		return result, nil
	}

	return nil, fmt.Errorf("project not found: %s", project)
}

// Reset clears all vault data
func (v *Vault) Reset() error {
	if !v.unlocked {
		return fmt.Errorf("vault is locked")
	}

	// Clear in-memory data
	v.data.Projects = make(map[string]map[string]string)

	// Save empty vault
	if err := v.save(); err != nil {
		return err
	}

	// Lock vault after reset
	v.Lock()
	return nil
}

// save saves the vault to disk
func (v *Vault) save() error {
	if !v.unlocked {
		return fmt.Errorf("vault is locked")
	}

	// Marshal vault data
	jsonData, err := json.Marshal(v.data)
	if err != nil {
		return fmt.Errorf("failed to marshal vault data: %w", err)
	}

	// Encrypt data
	encryptedData, err := crypto.Encrypt(jsonData, v.key)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	// Create encrypted vault structure
	encVault := EncryptedVault{
		Salt: v.data.Salt,
		Hash: v.data.Hash,
		Data: base64.StdEncoding.EncodeToString(encryptedData),
	}

	// Marshal encrypted vault
	vaultJSON, err := json.MarshalIndent(encVault, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted vault: %w", err)
	}

	// Write to file with proper permissions
	return os.WriteFile(v.path, vaultJSON, 0600)
}

// loadEncrypted loads the encrypted vault from disk
func (v *Vault) loadEncrypted() (*EncryptedVault, error) {
	data, err := os.ReadFile(v.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	var encVault EncryptedVault
	if err := json.Unmarshal(data, &encVault); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted vault: %w", err)
	}

	return &encVault, nil
}

// contains checks if str contains substr (case-insensitive)
func contains(str, substr string) bool {
	return len(substr) > 0 && len(str) >= len(substr) &&
		(str == substr || containsIgnoreCase(str, substr))
}

func containsIgnoreCase(str, substr string) bool {
	if len(substr) > len(str) {
		return false
	}

	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(str[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

// Exists checks if vault file exists
func (v *Vault) Exists() bool {
	_, err := os.Stat(v.path)
	return err == nil
}
