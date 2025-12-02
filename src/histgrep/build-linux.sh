#!/bin/bash

# This script builds the Go application for Linux.

set -e # Exit immediately if a command exits with a non-zero status.

# --- Configuration ---
VERSION=${APP_VERSION:-"0.5.0"}

if command -v git >/dev/null 2>&1; then
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    COMMIT=$(git rev-parse --short HEAD)
    VERSION="${VERSION}-${COMMIT}"
  fi
fi

OUTPUT_NAME="histgrep"
LDFLAGS_VAR="-s -w -X 'github.com/TJN25/histgrep/cmd.AppVersion=${VERSION}'"


# --- Build Process ---
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o "${OUTPUT_NAME}" -ldflags="${LDFLAGS_VAR}" .

# Make the final binary executable
chmod +x "${OUTPUT_NAME}"

echo ""
echo "✅ Build complete: ./${OUTPUT_NAME}"

