#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}🚀 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if version is provided
if [ $# -eq 0 ]; then
    print_error "Please provide a version number!"
    echo "Usage: $0 <version>"
    echo "Example: $0 1.0.7"
    exit 1
fi

NEW_VERSION=$1

# Validate version format (simple check)
if [[ ! $NEW_VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    print_error "Invalid version format! Use semantic versioning (e.g., 1.0.7)"
    exit 1
fi

print_status "Starting release process for version $NEW_VERSION"

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    print_warning "You are not on the main branch (current: $CURRENT_BRANCH)"
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Aborted by user"
        exit 1
    fi
fi

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    print_error "Working directory is not clean. Please commit or stash your changes."
    git status --short
    exit 1
fi

# Pull latest changes
print_status "Pulling latest changes from origin..."
git pull origin $CURRENT_BRANCH

# Get current version from package.json
CURRENT_VERSION=$(node -p "require('./package.json').version")
print_status "Current version: $CURRENT_VERSION"
print_status "New version: $NEW_VERSION"

# Check if tag already exists
if git rev-parse "v$NEW_VERSION" >/dev/null 2>&1; then
    print_error "Tag v$NEW_VERSION already exists!"
    exit 1
fi

# Generate and preview release notes
print_status "Generating release notes preview..."

# Get previous tag for changelog
PREV_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

# Create preview release notes
PREVIEW_FILE="/tmp/uzp-release-$NEW_VERSION.md"

cat << EOF > $PREVIEW_FILE
# UZP v$NEW_VERSION - User's Zecure Pocket

## 🆕 What's New

EOF

if [ -n "$PREV_TAG" ]; then
    echo "Changes since $PREV_TAG:" >> $PREVIEW_FILE
    echo "" >> $PREVIEW_FILE
    
    # Generate enhanced changelog with PR links
    git log $PREV_TAG..HEAD --oneline --no-merges | while read commit; do
        commit_hash=$(echo "$commit" | cut -d' ' -f1)
        commit_msg=$(echo "$commit" | cut -d' ' -f2-)
        
        # Try to extract PR number from commit message
        pr_num=$(echo "$commit_msg" | grep -o '#[0-9]\+' | head -1)
        
        if [ -n "$pr_num" ]; then
            pr_link="[${pr_num}](https://github.com/hungnguyen18/uzp-cli/pull/${pr_num#\#})"
            echo "- $commit_msg $pr_link" >> $PREVIEW_FILE
        else
            echo "- $commit_msg" >> $PREVIEW_FILE
        fi
    done
else
    echo "- Initial release" >> $PREVIEW_FILE
fi

cat << EOF >> $PREVIEW_FILE

## 📦 Installation

### NPM (Recommended)
\`\`\`bash
npm install -g uzp-cli
\`\`\`

### Direct Download
Download the appropriate binary for your platform from the assets below.

## 🚀 Quick Start
\`\`\`bash
uzp init
uzp add
uzp inject -p myapp > .env
\`\`\`

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

## 🔗 Links
- 📖 [Documentation](https://github.com/hungnguyen18/uzp-cli#readme)
- 🐛 [Report Issues](https://github.com/hungnguyen18/uzp-cli/issues)
- 💬 [Discussions](https://github.com/hungnguyen18/uzp-cli/discussions)
EOF

print_success "Release notes preview generated!"
echo ""
echo "📋 Release Notes Preview:"
echo "=========================="
cat $PREVIEW_FILE
echo "=========================="
echo ""

# Ask user if they want to edit
read -p "🤔 Do you want to edit the release notes? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Open editor (try different editors)
    if command -v code >/dev/null 2>&1; then
        print_status "Opening VS Code for editing..."
        code $PREVIEW_FILE
        read -p "Press Enter after you finish editing in VS Code..."
    elif command -v nano >/dev/null 2>&1; then
        print_status "Opening nano for editing..."
        nano $PREVIEW_FILE
    elif command -v vim >/dev/null 2>&1; then
        print_status "Opening vim for editing..."
        vim $PREVIEW_FILE
    else
        print_warning "No suitable editor found. You can manually edit: $PREVIEW_FILE"
        read -p "Press Enter when you're ready to continue..."
    fi
    
    print_success "Release notes updated!"
fi

# Confirm before proceeding
echo ""
print_status "Final release notes:"
echo "===================="
cat $PREVIEW_FILE
echo "===================="
echo ""

read -p "🚀 Proceed with release v$NEW_VERSION? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_error "Release cancelled by user"
    rm -f $PREVIEW_FILE
    exit 1
fi

# Update package.json version
print_status "Updating package.json version..."
if command -v jq >/dev/null 2>&1; then
    # Use jq if available
    jq ".version = \"$NEW_VERSION\"" package.json > package.json.tmp && mv package.json.tmp package.json
else
    # Fallback to sed
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s/\"version\": \"$CURRENT_VERSION\"/\"version\": \"$NEW_VERSION\"/" package.json
    else
        # Linux
        sed -i "s/\"version\": \"$CURRENT_VERSION\"/\"version\": \"$NEW_VERSION\"/" package.json
    fi
fi

# Verify the change
UPDATED_VERSION=$(node -p "require('./package.json').version")
if [ "$UPDATED_VERSION" != "$NEW_VERSION" ]; then
    print_error "Failed to update package.json version!"
    exit 1
fi

print_success "Updated package.json version to $NEW_VERSION"

# Commit the version change
print_status "Committing version change..."
git add package.json
git commit -m "chore: bump version to $NEW_VERSION"

# Push the commit
print_status "Pushing commit to origin..."
git push origin $CURRENT_BRANCH

# Create and push tag
print_status "Creating and pushing tag v$NEW_VERSION..."
git tag "v$NEW_VERSION"
git push origin "v$NEW_VERSION"

# Store custom release notes for GitHub Actions (if available)
if [ -f "$PREVIEW_FILE" ]; then
    RELEASE_NOTES_DIR=".github/release-notes"
    mkdir -p "$RELEASE_NOTES_DIR"
    cp "$PREVIEW_FILE" "$RELEASE_NOTES_DIR/v$NEW_VERSION.md"
    git add "$RELEASE_NOTES_DIR/v$NEW_VERSION.md"
    print_success "Custom release notes saved for GitHub Actions"
fi

print_success "🎉 Release process initiated!"
echo ""
echo "The GitHub Actions workflow will now:"
echo "  1. 🔨 Build cross-platform binaries"
echo "  2. 📦 Create GitHub release with your custom description"
echo "  3. 📤 Upload binaries to the release"
echo "  4. 🚀 Publish to NPM (reusing existing workflow)"
echo "  5. 🤖 Request Copilot analysis for enhanced release notes"
echo ""
echo "📋 Monitor progress at:"
echo "  🔗 GitHub Actions: https://github.com/hungnguyen18/uzp-cli/actions"
echo "  📦 Releases: https://github.com/hungnguyen18/uzp-cli/releases"
echo "  📡 NPM: https://www.npmjs.com/package/uzp-cli"
echo ""

# Cleanup
rm -f "$PREVIEW_FILE"

print_success "Done! 🚀" 