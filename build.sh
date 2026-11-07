#!/bin/bash
# Build script for Borg / HyperHarness / HyperCode
# Usage: ./build.sh [clean]

set -e

VERSION=$(cat VERSION 2>/dev/null || echo "0.5.0-alpha.5")
BINARY="hypercode"
BINARY_ALT="hyperharness"
LDFLAGS="-s -w -X github.com/robertpelloni/hyperharness/internal/buildinfo.Version=${VERSION}"

echo "╔════════════════════════════════════════════════════════════╗"
echo "║             Borg Build System v${VERSION}                  ║"
echo "╚════════════════════════════════════════════════════════════╝"

if [ "$1" = "clean" ]; then
	echo "Cleaning build cache..."
	go clean -cache
fi

echo ""
echo "🏗  Building ${BINARY} and ${BINARY_ALT}..."
go build -buildvcs=false -ldflags="${LDFLAGS}" -o ${BINARY} .
cp -f ${BINARY} ${BINARY_ALT} 2>/dev/null || true

if [ $? -eq 0 ]; then
	SIZE=$(ls -la ${BINARY} | awk '{print $5}')
	SIZE_MB=$(echo "scale=1; ${SIZE}/1048576" | bc 2>/dev/null || echo "?.?")
	echo "✅ Build successful: ${BINARY} (${SIZE_MB} MB)"
	echo ""
	echo "Run: ./${BINARY}"
	echo "   or: ./${BINARY} tui"
	echo "   or: ./${BINARY} serve"
else
	echo "❌ Build failed!"
	exit 1
fi
