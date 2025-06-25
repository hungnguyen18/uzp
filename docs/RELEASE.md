# Release Guide - UZP CLI

This document explains the fully automated release process for UZP CLI.

## 🚀 Quick Release (One Command)

```bash
# Release version 1.0.7
./scripts/release.sh 1.0.7

# Or using npm script
npm run release 1.0.7
```

**New Features:**
- 📝 **Interactive Preview** - See and edit release notes before publishing
- 🔗 **PR Links** - Automatic linking to Pull Requests in changelog
- 🤖 **Copilot Integration** - AI-powered release notes analysis
- ✅ **Confirmation Steps** - Review everything before proceeding

**That's it!** 🎉 Everything else is automated by GitHub Actions.

## 🎯 Interactive Release Workflow

### 📝 Preview & Edit Process

When you run the release script, you'll see:

```
🚀 Generating release notes preview...
✅ Release notes preview generated!

📋 Release Notes Preview:
==========================
# UZP v1.0.7 - User's Zecure Pocket

## 🆕 What's New

Changes since v1.0.6:

- feat: add vault export functionality [#123](https://github.com/hungnguyen18/uzp-cli/pull/123)
- fix: prevent clipboard memory leak [#124](https://github.com/hungnguyen18/uzp-cli/pull/124)
- docs: update installation guide

## 📦 Installation
[... full template ...]
==========================

🤔 Do you want to edit the release notes? (y/N): y
🚀 Opening VS Code for editing...
Press Enter after you finish editing in VS Code...
✅ Release notes updated!

🚀 Final release notes:
[... shows your edited version ...]

🚀 Proceed with release v1.0.7? (y/N): y
```

### 📝 Editing Options

The script will try these editors in order:
1. **VS Code** (`code`) - Opens in external editor
2. **Nano** (`nano`) - Terminal-based editor
3. **Vim** (`vim`) - Advanced terminal editor
4. **Manual** - Provides file path for manual editing

### 🔗 Automatic Features

- **PR Links**: Automatically detects `#123` in commit messages and links to PRs
- **Smart Changelog**: Groups commits since last tag
- **Template**: Includes installation, features, and platform support
- **Validation**: Confirms everything before proceeding

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
- 📦 **Release**: Creates GitHub release with your custom description (or auto-generated)
- 📤 **Upload**: Binaries attached to release assets
- 🚀 **Publish**: Reuses existing NPM workflow for reliability
- 🤖 **Copilot**: Requests AI analysis for enhanced release notes

## 🤖 GitHub Copilot Integration

### Automatic Analysis Request

After the release is created, GitHub Actions will automatically request GitHub Copilot to analyze the changes:

```markdown
@github-copilot Please analyze the changes for UZP v1.0.7 and categorize them:

## Changes to analyze:
Based on the commits and PRs included in this release, please provide:

### 🚀 New Features
List any new functionality added

### 🐛 Bug Fixes  
List any bugs that were fixed

### ⚡ Performance Improvements
List any performance optimizations

### 🔒 Security Updates
List any security-related changes

### 💥 Breaking Changes
List any breaking changes (if any)

### 📚 Documentation
List documentation updates

### 🧹 Other Changes
List other minor changes, refactoring, etc.

Please format this as markdown and focus on user-facing changes.
```

### How to Use Copilot Analysis

1. **After release**: Go to the GitHub release page
2. **Check comments**: Look for Copilot's analysis in the discussion
3. **Copy insights**: Use Copilot's categorization for announcements
4. **Future releases**: Reference previous analyses for consistency

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

## 🎊 What's New in Release v2.0

We've completely redesigned the release process with these major improvements:

### ✨ **New Features**

| Feature | Description |
|---------|-------------|
| 📝 **Interactive Preview** | See and edit release notes before publishing |
| 🔗 **Auto PR Links** | Automatic linking to Pull Requests in changelog |
| 🤖 **Copilot Integration** | AI-powered release notes analysis and categorization |
| ✅ **Smart Confirmation** | Multiple validation steps before release |
| 🔄 **Workflow Reuse** | NPM publishing reuses existing reliable workflow |
| 🎨 **Multi-Editor Support** | VS Code, Nano, Vim, or manual editing |

### 🛡️ **Improved Reliability**

- **Zero-downtime releases**: Reuses battle-tested publish workflow
- **Smart validation**: Checks version, authorization, and git status
- **Safe rollback**: Clear instructions for emergency rollback
- **Better error handling**: More descriptive error messages

### 🚀 **Enhanced Developer Experience**

- **One command release**: `./scripts/release.sh 1.0.7`
- **Real-time preview**: See exactly what will be published
- **Easy customization**: Edit release notes before publishing
- **Progress tracking**: Clear workflow status and links

---

## 📋 Quick Reference

```bash
# 🚀 Start release
./scripts/release.sh 1.0.8

# 📝 Script will show preview
# 🤔 Choose to edit or proceed
# ✅ Confirm final version
# 🎉 Automatic GitHub Actions takes over!

# 📊 Monitor progress
# - GitHub Actions: /actions
# - Releases: /releases  
# - NPM: npmjs.com/package/uzp-cli
```

---

**Questions about releases?** Open an issue or discussion on GitHub! 🚀 