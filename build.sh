#!/bin/bash

# M3TAL Platform Build Script (Linux/WSL)

# Try to find go if not in path
GO_CMD=$(which go 2>/dev/null)
if [ -z "$GO_CMD" ]; then
    echo "❌ Go not found in PATH."
    echo "Try: export PATH=\$PATH:/usr/local/go/bin"
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
