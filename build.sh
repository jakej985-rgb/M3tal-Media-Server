#!/bin/bash

# M3TAL Platform Build Script (Linux/WSL)

# Try to find go if not in path
GO_CMD=$(which go 2>/dev/null)
if [ -z "$GO_CMD" ]; then
    echo "❌ Go not found in PATH."
    read -p "Would you like to attempt to install Go now? (y/n): " confirm
    if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
        if [ -x "$(command -v apt)" ]; then
            echo "📦 Detected apt. Installing golang..."
            sudo apt update && sudo apt install -y golang
            GO_CMD=$(which go 2>/dev/null)
        else
            echo "⚠️  Package manager not supported for auto-install."
            echo "Please install Go manually: https://go.dev/doc/install"
            exit 1
        fi
    else
        echo "Try: export PATH=\$PATH:/usr/local/go/bin"
        exit 1
    fi
fi

if [ -z "$GO_CMD" ]; then
    echo "❌ Go installation failed or not found in PATH."
    exit 1
fi

echo "🚀 Building M3TAL CLI..."
$GO_CMD build -o m3tal ./cmd/m3tal

echo "🚀 Building M3TAL API..."
$GO_CMD build -o m3tal-api ./cmd/api

if [ $? -eq 0 ]; then
    echo "✅ Build complete."
else
    echo "❌ Build failed."
    exit 1
fi
