#!/bin/bash

set -e

VERSION=${1:-"1.0.0"}
BINARY_NAME="uzp"
BUILD_DIR="build"

echo "🔨 Building UZP v$VERSION for multiple platforms..."

# Clean build directory
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# Platforms to build for
PLATFORMS=(
  "linux/amd64"
  "linux/arm64"  
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
  "windows/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
  IFS='/' read -r GOOS GOARCH <<< "$PLATFORM"
  
  OUTPUT_NAME="$BINARY_NAME-$GOOS-$GOARCH"
  if [ "$GOOS" = "windows" ]; then
    OUTPUT_NAME="$OUTPUT_NAME.exe"
  fi
  
  echo "📦 Building $OUTPUT_NAME..."
  
  env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$BUILD_DIR/$OUTPUT_NAME" .
  
  echo "✅ Built $OUTPUT_NAME"
done

echo ""
echo "🎉 Build complete! Binaries available in $BUILD_DIR/"
echo ""
echo "📋 Built binaries:"
ls -la $BUILD_DIR/

echo ""
echo "🚀 To create a GitHub release:"
echo "   1. Create a new tag: git tag v$VERSION"
echo "   2. Push tag: git push origin v$VERSION"  
echo "   3. Create GitHub release and upload binaries from $BUILD_DIR/" 