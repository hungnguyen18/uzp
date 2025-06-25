# UZP-CLI - User's Zecure Pocket
[![npm version](https://badge.fury.io/js/uzp-cli.svg)](https://badge.fury.io/js/uzp-cli)
[![npm downloads](https://img.shields.io/npm/dm/uzp-cli.svg)](https://www.npmjs.com/package/uzp-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hungnguyen18/uzp-cli)](https://golang.org/)
[![Go 1.23.10+](https://img.shields.io/badge/Go-1.23.10+-blue.svg)](https://golang.org/dl/)
[![Security](https://img.shields.io/badge/Security-AES--256--GCM-green.svg)](SECURITY.md)
[![Contributing](https://img.shields.io/badge/Contributing-Welcome-brightgreen.svg)](CONTRIBUTING.md)

A professional command-line tool for securely storing and managing sensitive information such as API keys, access tokens, and service credentials. All data is encrypted using AES-256-GCM and stored locally.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Security](#security)
- [Examples](#examples)
- [Contributing](#contributing)
- [Releases](#releases)
- [Support](#support)

## Features

### Core Capabilities
- 🔐 **AES-256-GCM encryption** with scrypt key derivation (N=32768, r=8, p=1)
- 🔑 **Master password protection** - never stored, only hashed
- 🔄 **Auto-unlock** - prompts for password when needed, no session management
- 📁 **Project-based organization** - group secrets by application/service
- 📋 **Clipboard integration** with automatic clearing after TTL
- 🔍 **Search functionality** for quick access across all projects
- 📄 **Environment file export** (.env generation) for development workflows
- 🌍 **Cross-platform support** (macOS, Linux, Windows)
- 🔒 **Secure file permissions** - vault files created with 0600 permissions

### Security Features
- **Memory safety**: Sensitive data cleared from memory immediately after use
- **No password storage**: Only password hash stored for verification
- **Secure vault location**: 
  - macOS/Linux: `~/.uzp/uzp.vault`
  - Windows: `%USERPROFILE%\.uzp\uzp.vault`

## Installation

```bash
# NPM (Recommended)
npm install -g uzp-cli

# Alternative package managers
yarn global add uzp-cli    # Yarn
pnpm add -g uzp-cli        # PNPM  
bun add -g uzp-cli         # Bun

# Manual installation
git clone https://github.com/hungnguyen18/uzp-cli.git
cd uzp-cli && go build -o uzp

# NPM Registry Configuration (if needed)
cp .npmrc.example .npmrc   # Edit for custom registries
```

## Quick Start

```bash
# 1. Initialize vault with master password
uzp init

# 2. Add your first secret
uzp add
# Project: myapp
# Key: api_key 
# Value: sk-1234567890abcdef

# 3. Use your secrets
uzp get myapp/api_key              # Display secret
uzp copy myapp/api_key             # Copy to clipboard
uzp inject -p myapp > .env         # Export as .env file
```

## Commands

| Command | Description | Example |
|---------|-------------|---------|
| `uzp init` | Initialize new vault | `uzp init` |
| `uzp add` | Add a secret | `uzp add` |
| `uzp get <project/key>` | Get secret value | `uzp get myapp/api_key` |
| `uzp copy <project/key>` | Copy to clipboard | `uzp copy myapp/api_key` |
| `uzp update <project/key>` | Update existing secret | `uzp update myapp/api_key` |
| `uzp list` | List all secrets | `uzp list` |
| `uzp search <keyword>` | Search secrets | `uzp search api` |
| `uzp inject -p <project>` | Export to .env format | `uzp inject -p myapp > .env` |
| `uzp reset` | Delete all data | `uzp reset` |

## Security

UZP-CLI follows security-first principles:

- **🔐 Encryption**: AES-256-GCM with random salts and nonces
- **🔑 Key Derivation**: scrypt with secure parameters (N=32768, r=8, p=1)  
- **🛡️ Password Protection**: Master password never stored, only its hash
- **🧹 Memory Safety**: Sensitive data cleared from memory after use
- **📁 File Permissions**: Vault files created with 0600 (user-only access)
- **📋 Clipboard Safety**: Automatic clearing after configurable TTL

### Security Warnings

- ⚠️ **Never share your master password**
- 🔒 **Keep your vault file secure and backed up**
- 🔑 **Use a strong, unique master password (12+ characters recommended)**
- 🚫 **Don't store your master password in scripts or files**

For security issues, see our [Security Policy](SECURITY.md).

## Examples

### Basic Workflow
```bash
# Initialize and add secrets
uzp init
uzp add  # myapp/api_key
uzp add  # myapp/database_url
uzp add  # aws/access_key_id

# Use secrets in development
uzp inject -p myapp > .env.local
uzp inject -p aws > aws.env
uzp copy myapp/api_key

# Search and manage
uzp list                    # View all secrets
uzp search database         # Find specific secrets
uzp update myapp/api_key    # Update existing values
```

### Environment File Export
```bash
# Export project secrets
uzp inject -p myapp > .env

# Multiple environments
uzp inject -p myapp > .env.local
uzp inject -p myapp-prod > .env.production

# Preview before export
uzp inject -p myapp
```

**Generated .env format:**
```bash
# Environment variables for project: myapp
# Generated by uzp
API_KEY=your_secret_value
DATABASE_URL=postgresql://user:pass@host:5432/db
```

## Contributing

**New contributors:** Get started quickly with the [**Quick Start section**](CONTRIBUTING.md#-quick-start---get-contributing-in-5-minutes) in our Contributing Guide! ⚡

Our [Contributing Guide](CONTRIBUTING.md) covers everything from 5-minute setup to comprehensive development practices, security requirements, and submission process.

**Thank you for helping make UZP-CLI more secure! 🔐**

## Releases

**Release Information:**
- 🔔 **Latest**: Check [GitHub Releases](https://github.com/hungnguyen18/uzp-cli/releases) for newest version
- 📅 **Schedule**: Monthly minor releases, patches as needed for critical bugs
- 📦 **Versioning**: Follows [Semantic Versioning](https://semver.org/) (vMAJOR.MINOR.PATCH)
- 📝 **Notes**: Detailed release notes with features, fixes, and contributor credits

```bash
# Check your installed version
uzp --version

# Update to latest version  
npm update -g uzp-cli
```

## Support

**Get Help:**
- 🐛 [Bug Reports](https://github.com/hungnguyen18/uzp-cli/issues/new) - Report issues
- 💡 [Feature Requests](https://github.com/hungnguyen18/uzp-cli/issues) - Suggest improvements  
- ❓ [Questions](https://github.com/hungnguyen18/uzp-cli/discussions) - Ask the community
- 🔒 [Security Issues](SECURITY.md) - Private security reporting

**Resources:**
- 📖 [Contributing Guidelines](CONTRIBUTING.md) - Development and contribution guide
- 🔐 [Security Policy](SECURITY.md) - Security practices and vulnerability reporting
- 📦 [NPM Package](https://www.npmjs.com/package/uzp-cli) - Official package
- 🏗️ [Technical Docs](docs/) - Internal documentation for maintainers
- 📜 [License](LICENSE) - MIT License

---

**UZP-CLI** - Your secrets, secured locally. 🔐 