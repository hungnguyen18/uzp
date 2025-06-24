# UZP - Secure CLI Tool for Managing Secrets

`uzp` is a command-line tool designed to securely store and manage sensitive information such as API keys, access tokens, and service credentials. All data is encrypted using AES-256-GCM and stored locally.

## Features

- 🔐 **AES-256-GCM encryption** for maximum security
- 🔑 **Master password protection** with scrypt key derivation
- 📋 **Clipboard integration** with automatic clearing
- 🔍 **Search functionality** for quick access
- 📁 **Project-based organization** of secrets
- 🌍 **Cross-platform support** (macOS, Linux, Windows)

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from source

```bash
git clone https://github.com/hungnguyen/uzp.git
cd uzp
go build -o uzp
```

### Install globally

```bash
go install .
```

## Usage

### Initialize a new vault

```bash
uzp init
```

This creates a new encrypted vault with your master password.

### Unlock the vault

```bash
uzp unlock
```

Enter your master password to unlock the vault for the current session.

### Add a secret

```bash
uzp add
```

You'll be prompted for:
- Project name (e.g., "myapp", "aws", "github")
- Key name (e.g., "api_key", "access_token")
- Value (hidden input)

### Get a secret

```bash
uzp get myapp/api_key
```

### Copy to clipboard

```bash
uzp copy myapp/api_key
```

The value will be copied to clipboard and automatically cleared after 15 seconds (configurable with `--ttl`).

### List all secrets

```bash
uzp list
```

### Search for secrets

```bash
uzp search api
```

### Export project secrets as .env

```bash
uzp inject --project myapp > .env
```

### Lock the vault

```bash
uzp lock
```

### Reset vault (danger!)

```bash
uzp reset
```

This will delete ALL secrets permanently after confirmation.

## Security

- **Encryption**: AES-256-GCM for data encryption
- **Key Derivation**: scrypt with secure parameters (N=32768, r=8, p=1)
- **Password Protection**: Master password never stored, only its hash
- **Clipboard Safety**: Automatic clearing after TTL
- **File Permissions**: Vault file created with 0600 permissions

## Storage Location

The vault is stored at:
- macOS/Linux: `~/.uzp/uzp.vault`
- Windows: `%USERPROFILE%\.uzp\uzp.vault`

## Development

### Project Structure

```
uzp/
├── cmd/               # CLI commands
│   ├── root.go
│   ├── init.go
│   ├── unlock.go
│   └── ...
├── internal/
│   ├── crypto/        # Encryption/decryption
│   ├── storage/       # Vault management
│   └── utils/         # Utilities
├── main.go
└── go.mod
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o uzp
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Security Notice

- Never share your master password
- Keep your vault file (`~/.uzp/uzp.vault`) secure
- Regularly backup your vault file
- Use a strong, unique master password 