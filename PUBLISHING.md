# Publishing Guide

## Setup for Publishing

### 1. GitHub Repository Setup

1. Create GitHub repository: `https://github.com/hungnguyen18/uzp`
2. Push code to GitHub:
```bash
git add .
git commit -m "Initial release v1.0.0"
git push origin main
```

### 2. Build Release Binaries

```bash
# Build all platform binaries
./scripts/build.sh 1.0.0

# Check built binaries
ls -la build/
```

### 3. Create GitHub Release

1. Create and push tag:
```bash
git tag v1.0.0
git push origin v1.0.0
```

2. Go to GitHub releases: `https://github.com/hungnguyen18/uzp/releases`
3. Click "Create a new release"
4. Select tag: `v1.0.0`
5. Release title: `UZP v1.0.0 - User's Zecure Pocket`
6. Description:
```markdown
## 🚀 UZP v1.0.0 - First Release

User's Zecure Pocket - A secure CLI tool for managing secrets.

### Features
- 🔐 AES-256-GCM encryption
- 🔄 Auto-unlock workflow
- 📄 Environment file export (.env)
- 📋 Clipboard integration
- 🔍 Search functionality

### Installation
```bash
npm install -g @hungnguyen18/uzp
```

### Quick Start
```bash
uzp init
uzp add
uzp inject -p myapp > .env
```

### Platform Support
- macOS (Intel & Apple Silicon)
- Linux (x64 & ARM64)
- Windows (x64 & ARM64)
```

7. Upload binaries from `build/` directory:
   - `uzp-darwin-amd64`
   - `uzp-darwin-arm64`
   - `uzp-linux-amd64`
   - `uzp-linux-arm64`
   - `uzp-windows-amd64.exe`
   - `uzp-windows-arm64.exe`

8. Click "Publish release"

### 4. Publish to NPM

1. Login to npm:
```bash
npm login
```

2. Publish package:
```bash
npm publish
```

3. Verify installation:
```bash
npm install -g @hungnguyen18/uzp
uzp --help
```

## Release Workflow

For future releases:

1. Update version in `package.json`
2. Build binaries: `./scripts/build.sh <version>`
3. Create GitHub release with new binaries
4. Publish to npm: `npm publish`

## Testing

### Test Local Package
```bash
# Pack locally
npm pack

# Install from local package
npm install -g ./uzp-1.0.0.tgz

# Test commands
uzp --help
uzp init
```

### Test NPM Installation
```bash
# After publishing
npm install -g @hungnguyen18/uzp
uzp --help
```

## Troubleshooting

### NPM Installation Fails
- Check GitHub release exists with correct binary names
- Verify binary naming matches `uzp-<platform>-<arch>[.exe]`
- Check network connectivity

### Binary Not Found
- Ensure binaries are uploaded to GitHub release
- Verify platform detection in `scripts/install.js`
- Check file permissions (should be 755)

## Package Structure

```
hungnguyen18-uzp-1.0.0.tgz
├── package.json
├── README.md
├── LICENSE
├── bin/uzp (placeholder)
├── scripts/
│   ├── install.js (downloads binary)
│   └── uninstall.js (cleanup)
└── .npmignore
``` 