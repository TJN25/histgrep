#!/bin/bash

set -e 

# --- Configuration ---
VERSION=${APP_VERSION:-"0.5.0"}

if command -v git >/dev/null 2>&1; then
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    COMMIT=$(git rev-parse --short HEAD)
    VERSION="${VERSION}-${COMMIT}"
  fi
fi

OUTPUT_NAME="histgrep"
TMP_AMD64_NAME="${OUTPUT_NAME}_amd64"
TMP_ARM64_NAME="${OUTPUT_NAME}_arm64"

# The package path to the AppVersion variable
LDFLAGS_VAR="-s -w -X 'github.com/TJN25/histgrep/cmd.AppVersion=${VERSION}'"

# --- Build Process ---
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o "${TMP_AMD64_NAME}" -ldflags="${LDFLAGS_VAR}" .

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o "${TMP_ARM64_NAME}" -ldflags="${LDFLAGS_VAR}" .

echo "Creating universal binary with lipo..."
lipo -create -output "${OUTPUT_NAME}" "${TMP_AMD64_NAME}" "${TMP_ARM64_NAME}"

echo "Cleaning up temporary files..."
rm "${TMP_AMD64_NAME}" "${TMP_ARM64_NAME}"

# Make the final binary executable
chmod +x "${OUTPUT_NAME}"

echo ""
echo "✅ Universal build complete: ./${OUTPUT_NAME}"
echo "Run 'file ./${OUTPUT_NAME}' to verify the architectures."

