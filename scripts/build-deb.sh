#!/bin/bash
set -e

echo "🚀 Building M3TAL Core .deb with nfpm..."

# 1. Compile binaries
go build -o m3tal ./cmd/m3tal
go build -o m3tal-api ./cmd/api

# 2. Get version
VERSION=$(cat VERSION)
export VERSION

# 3. Build the package using nfpm
# We use -f packaging/nfpm.yaml to ensure it uses our refined config
if command -v nfpm >/dev/null 2>&1; then
    nfpm pkg --packager deb --target "m3tal_${VERSION}_amd64.deb" -f packaging/nfpm.yaml
else
    echo "❌ nfpm not found. Please install it: https://nfpm.goreleaser.com/install/"
    exit 1
fi

echo "✅ Package built: m3tal_${VERSION}_amd64.deb"
