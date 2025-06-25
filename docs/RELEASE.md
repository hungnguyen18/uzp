# Release Guide - UZP CLI

This document explains the fully automated release process for UZP CLI.

## 🚀 Quick Release (One Command)

```bash
# Release version 1.0.7
./scripts/release.sh 1.0.7

# Or using npm script
npm run release 1.0.7
```

**That's it!** 🎉 Everything else is automated by GitHub Actions.

## 🔄 What Happens Automatically

### 1. **Local Script** (`./scripts/release.sh`)
- ✅ Validates version format
- ✅ Checks if you're on main branch  
- ✅ Ensures working directory is clean
- ✅ Updates `package.json` version
- ✅ Commits version change
- ✅ Creates and pushes git tag `v1.0.7`

### 2. **GitHub Actions** (`.github/workflows/auto-release.yml`)
- 🔨 **Build**: Cross-platform binaries for all platforms
- 📦 **Release**: Creates GitHub release with auto-generated description
- 📤 **Upload**: Binaries attached to release assets
- 🚀 **Publish**: Package published to NPM automatically

## 📋 Prerequisites

### Required Setup
- [ ] NPM token in GitHub Secrets (`NPM_TOKEN`)
- [ ] Authorized user (`hungnguyen18` only)
- [ ] Clean working directory
- [ ] On `main` branch (recommended)

### Required Tools
- `git` - For version control
- `node` & `npm` - For package.json operations
- `jq` (optional) - For JSON manipulation (fallback to `sed`)

## 📝 Release Template

The GitHub release is automatically created with this template:

```markdown
## 🆕 What's New
- [Auto-generated from git commits since last tag]

## 📦 Installation

### NPM (Recommended)
```bash
npm install -g uzp-cli
```

### Direct Download
Download the appropriate binary for your platform from the assets below.

## 🚀 Quick Start
```bash
uzp init
uzp add
uzp inject -p myapp > .env
```

## ✨ Features
- 🔐 AES-256-GCM encryption
- 🔄 Auto-unlock workflow
- 📄 Environment file export (.env)
- 📋 Clipboard integration
- 🔍 Search functionality

## 🖥️ Platform Support
- **macOS**: Intel (x64) & Apple Silicon (ARM64)
- **Linux**: x64 & ARM64
- **Windows**: x64 & ARM64
```

## 🛠️ Supported Platforms

The release builds for all these platforms automatically:

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| **macOS** | Intel (x64) | `uzp-darwin-amd64` |
| **macOS** | Apple Silicon (ARM64) | `uzp-darwin-arm64` |
| **Linux** | x64 | `uzp-linux-amd64` |
| **Linux** | ARM64 | `uzp-linux-arm64` |
| **Windows** | x64 | `uzp-windows-amd64.exe` |
| **Windows** | ARM64 | `uzp-windows-arm64.exe` |

## 🔍 Monitor Release Progress

After running the release script, monitor progress at:

- 🔗 **GitHub Actions**: https://github.com/hungnguyen18/uzp-cli/actions
- 📦 **Releases**: https://github.com/hungnguyen18/uzp-cli/releases  
- 📡 **NPM Package**: https://www.npmjs.com/package/uzp-cli

## ⚠️ Important Notes

### Security Requirements
- Only `hungnguyen18` can trigger releases (GitHub authorization check)
- All security files require manual review before merge
- NPM token must be configured in GitHub Secrets

### Version Requirements
- Must use semantic versioning (e.g., `1.0.7`)
- Package.json version must match git tag version
- Tag must not already exist

### Quality Gates
- All CI tests must pass before NPM publish
- Build must succeed for all platforms
- No duplicate versions allowed on NPM

## 🛡️ Error Handling & Rollback

### Common Issues

**❌ "Tag already exists"**
```bash
# Delete local and remote tag
git tag -d v1.0.7
git push origin :refs/tags/v1.0.7
```

**❌ "Version mismatch"**
- Ensure `package.json` version matches tag version
- Run the release script again with correct version

**❌ "NPM publish failed"**
- Check NPM token is valid in GitHub Secrets
- Verify package name is available
- Check network connectivity

### Emergency Rollback

If you need to rollback a release:

```bash
# 1. Unpublish from NPM (within 72 hours)
npm unpublish uzp-cli@1.0.7

# 2. Delete GitHub release
# Go to GitHub → Releases → Delete release

# 3. Delete git tag
git tag -d v1.0.7
git push origin :refs/tags/v1.0.7

# 4. Revert version in package.json
git revert HEAD  # Reverts the version bump commit
```

## 📊 Release Analytics

After release, you can track:

- **NPM Downloads**: https://npm-stat.com/charts.html?package=uzp-cli
- **GitHub Downloads**: Check release assets download count
- **Version Usage**: `npm view uzp-cli versions --json`

## 🎯 Best Practices

### Before Release
- [ ] Test locally: `go test ./...`
- [ ] Update documentation if needed
- [ ] Review changelog/commits since last release
- [ ] Ensure no breaking changes (for patch/minor releases)

### After Release
- [ ] Verify NPM installation: `npm install -g uzp-cli@1.0.7`
- [ ] Test binary downloads work
- [ ] Update any dependent projects
- [ ] Announce release in relevant channels

### Release Cadence
- **Patch releases** (1.0.x): Bug fixes, security updates → As needed
- **Minor releases** (1.x.0): New features, improvements → Monthly
- **Major releases** (x.0.0): Breaking changes → When necessary

---

**Questions about releases?** Open an issue or discussion on GitHub! 🚀 