# Contributing to UZP-CLI

Welcome! 👋 Thank you for your interest in contributing to UZP-CLI. This guide will help you get started quickly and contribute effectively.

> 🆘 **Need help?** Check [❓ Getting Help](#-getting-help) section or our [GitHub Discussions](https://github.com/hungnguyen18/uzp-cli/discussions)

## 📖 Table of Contents

- [🚀 Quick Start (5 minutes)](#-quick-start-5-minutes)
- [📋 What Can You Contribute?](#-what-can-you-contribute)
- [🔍 Review Process](#-review-process)
- [🛠️ Development Guidelines](#️-development-guidelines)
- [📝 Commit Messages](#-commit-messages)
- [🔒 Security Guidelines](#-security-guidelines)
- [🚀 Release Process](#-release-process)
- [❓ Getting Help](#-getting-help)
- [🎯 Tips for New Contributors](#-tips-for-new-contributors)

---

## 🚀 Quick Start (5 minutes)

```bash
# 1. Fork on GitHub, then clone
git clone https://github.com/YOUR_USERNAME/uzp-cli.git
cd uzp-cli

# 2. Set up & test
go mod download && npm install
go build -o uzp . && ./uzp --help

# 3. Create your feature
git checkout -b feature/your-feature-name
# ... make your changes ...

# 4. Submit
git add . && git commit -m "feat: your description"
git push origin feature/your-feature-name
# Then create PR on GitHub
```

**That's it!** 🎉 Our CI will test everything automatically.

> 💡 **New to contributing?** Start with issues labeled [`good first issue`](https://github.com/hungnguyen18/uzp-cli/labels/good%20first%20issue) or improve documentation!

---

## 📋 What Can You Contribute?

**Easy to get started** ✅ (Auto-merge after CI):
- 🐛 **Bug Fixes** - Fix commands, error handling
- 📚 **Documentation** - README, help text, examples  
- 🧹 **Code Cleanup** - Refactoring, formatting
- ✨ **New Features** - New commands, utilities

**Needs extra review** 🔍 (Security files):
- 🔒 **Security/Crypto** - Encryption, storage core

> 💡 **New contributors:** Start with documentation or bug fixes! They're automatically merged when CI passes.

### Which files need manual review?
Only these security-critical files require owner review:
- `internal/crypto/` & `internal/storage/` (encryption/vault)
- `.github/workflows/` & `go.mod` (CI/dependencies)

---

## 🔍 Review Process

**Most PRs (90%):** Submit → CI Tests → ✅ Auto-merge  
**Security PRs (10%):** Submit → CI Tests → ✅ Owner Review → Merge

That's it! Our automated CI will test your code and most changes get merged automatically once tests pass. Only security-sensitive files need a human to review them.

---

## 🛠️ Development Guidelines

### Prerequisites
- **Go**: 1.23.10+ (required for security)
- **Node.js**: 18+
- **Git**: 2.30+

### Branch Naming
Use this format: `<type>/<description_in_snake_case>`

```bash
feature/add_backup_export     # New features
bug/fix_clipboard_leak        # Bug fixes
docs/update_readme           # Documentation
```

### Code Style

#### Go Code
```go
// ✅ Good: Clear names and error handling
func EncryptVaultData(data []byte, password string) ([]byte, error) {
    if len(data) == 0 {
        return nil, errors.New("data cannot be empty")
    }
    
    encrypted, err := encrypt(data, password)
    if err != nil {
        return nil, fmt.Errorf("encryption failed: %w", err)
    }
    
    return encrypted, nil
}

// ❌ Bad: Unclear names, poor error handling
func Process(d interface{}) interface{} {
    // ...
}
```

#### Security Requirements
```go
// ✅ Always clear sensitive data
password := getPassword()
defer func() {
    for i := range password {
        password[i] = 0
    }
}()

// ✅ Validate all inputs
func ValidateProjectName(name string) error {
    if len(name) == 0 {
        return errors.New("project name cannot be empty")
    }
    if len(name) > 255 {
        return errors.New("project name too long")
    }
    return nil
}
```

### Testing
```bash
# Run tests before submitting PR
go test ./...
npm test

# Test your changes manually
./uzp init
./uzp add test-entry
./uzp get test-entry
```

---

## 📝 Commit Messages

Use this format:
```
<type>: <description>

Examples:
feat: add vault export functionality
fix: prevent clipboard memory leak
docs: update installation guide
security: strengthen password validation
```

**Types:**
- `feat` - New features
- `fix` - Bug fixes
- `docs` - Documentation
- `security` - Security improvements
- `test` - Tests
- `refactor` - Code cleanup

---

## 🔒 Security Guidelines

UZP-CLI handles sensitive data, so security is critical:

### 1. Sensitive Data
- Always clear passwords/keys from memory after use
- Use secure random generation (`crypto/rand`)
- Validate all user inputs

### 2. File Operations
- Use secure file permissions (0600 for vault files)
- Validate file paths to prevent directory traversal
- Handle file errors properly

### 3. Dependencies
- Keep dependencies up to date
- Audit for vulnerabilities: `go mod audit`, `npm audit`

---

## 🚀 Release Process

### When Your Contribution Gets Released

| Type | Release Timeline |
|------|------------------|
| 🐛 **Critical Bug** | Days (patch release) |
| 🔒 **Security Fix** | Immediate |
| ✨ **New Feature** | Monthly (minor release) |
| 📚 **Documentation** | Next release |

### Versioning
- `v1.0.1` - Patch (bug fixes)
- `v1.1.0` - Minor (new features)
- `v2.0.0` - Major (breaking changes)

---

## ❓ Getting Help

- 📖 **Documentation**: Check README.md first
- 🐛 **Issues**: Search existing issues before creating new ones
- 💬 **Questions**: Use GitHub Discussions
- 🔒 **Security**: Email maintainers privately for security issues

---

## 📁 Project Structure

```
uzp-cli/
├── cmd/           # CLI commands (add.go, get.go, etc.)
├── internal/      # Core application code
│   ├── crypto/    # 🔒 Encryption (needs owner review)
│   ├── storage/   # 🔒 Vault storage (needs owner review)
│   └── utils/     # ✅ Utilities (auto-merge after CI)
├── scripts/       # ✅ Build scripts (auto-merge after CI)
└── docs/          # ✅ Documentation (auto-merge after CI)
```

---

## 🎯 Tips for New Contributors

**🌱 First time contributing to open source?**
1. **Start with docs** - Fix typos, improve examples, add missing info
2. **Look for `good first issue`** - These are designed for newcomers
3. **Read the code** - Browse `cmd/` folder to understand how commands work
4. **Ask questions** - Use GitHub Discussions or comment on issues

**⚡ Quick wins:**
- Fix a typo in README or help text
- Add an example to documentation  
- Improve error messages
- Add tests for existing functions

**🧪 Before submitting:**
- Test your changes: `go run . --help`
- Run tests: `go test ./...`
- Check your commit message follows the format

**❓ Not sure what to work on?** Check out our [good first issues](https://github.com/hungnguyen18/uzp-cli/labels/good%20first%20issue)!

---

---

## 🎉 Ready to Contribute?

**Thank you for contributing to UZP-CLI!** 🔐

Remember:
- 📚 **Start small** - Documentation and bug fixes are great first contributions
- 🤖 **CI does the work** - Most PRs merge automatically after tests pass  
- 💬 **Ask questions** - We're here to help you succeed
- 🔒 **Security matters** - But don't let it intimidate you!

Every contribution, big or small, helps make UZP-CLI better and more secure for everyone.

---

*For detailed technical guidelines, see our [docs/](docs/) directory.*
