#!/bin/bash
set -e

echo "🚀 Building M3TAL Core .deb..."

# 1. Clean build artifacts
BUILD_DIR="packaging/deb/build"
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR/DEBIAN"

# 2. Compile binaries
go build -o m3tal ./cmd/m3tal
go build -o m3tal-api ./cmd/api

# 3. Copy files to build directory structure
mkdir -p "$BUILD_DIR/usr/bin"
cp m3tal "$BUILD_DIR/usr/bin/"
cp m3tal-api "$BUILD_DIR/usr/bin/"

mkdir -p "$BUILD_DIR/lib/systemd/system"
cp packaging/m3tal.service "$BUILD_DIR/lib/systemd/system/"
cp packaging/m3tal-api.service "$BUILD_DIR/lib/systemd/system/"

mkdir -p "$BUILD_DIR/opt/m3tal/stack"
cp -r deploy/stack/* "$BUILD_DIR/opt/m3tal/stack/"

mkdir -p "$BUILD_DIR/usr/share/m3tal/defaults"
cp packaging/config.yaml "$BUILD_DIR/usr/share/m3tal/defaults/"

# 4. Copy control files
cp packaging/deb/DEBIAN/* "$BUILD_DIR/DEBIAN/"
chmod +x "$BUILD_DIR/DEBIAN/postinst"
chmod +x "$BUILD_DIR/DEBIAN/prerm"

# 5. Build the package
VERSION=$(cat VERSION)
dpkg-deb --build "$BUILD_DIR" "m3tal-core_${VERSION}_amd64.deb"

echo "✅ Package built: m3tal-core_${VERSION}_amd64.deb"
