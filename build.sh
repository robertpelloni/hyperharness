#!/bin/bash
# Build script for HyperHarness
# Usage: ./build.sh [clean]

set -e

VERSION=$(cat VERSION 2>/dev/null || echo "0.5.0-alpha.1")
BINARY="hyperharness.exe"
LDFLAGS="-s -w -X internal/buildinfo.Version=${VERSION}"

echo "╔════════════════════════════════════════════════════════════╗"
echo "║           HyperHarness Build System v${VERSION}           ║"
echo "╚════════════════════════════════════════════════════════════╝"

if [ "$1" = "clean" ]; then
	echo "Cleaning..."
	rm -f ${BINARY}
	go clean -cache
fi

echo ""
echo "🏗  Building..."
go build -buildvcs=false -ldflags="${LDFLAGS}" -o ${BINARY} .

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
