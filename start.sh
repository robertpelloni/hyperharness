#!/bin/bash
set -e

VERSION=$(cat VERSION 2>/dev/null || echo "0.5.0-alpha.5")
echo "Starting Borg Control Plane v${VERSION} (Go Native Core)..."

if ! command -v go &> /dev/null
then
    echo "go could not be found. Please install Go 1.22+."
    exit 1
fi

echo "Building Borg Go Control Plane..."
./build.sh

echo "Launching Borg..."
./hypercode "$@"
