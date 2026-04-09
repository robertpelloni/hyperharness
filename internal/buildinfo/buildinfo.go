// Package buildinfo provides version and build metadata.
// Set at build time via -ldflags "-X internal/buildinfo.Version=...".
//
// WHAT: Single source of truth for version number
// WHY: Version is referenced by CLI, HTTP API, RPC, and tools
// HOW: Variable set by linker flags, falls back to source default
package buildinfo

import (
	"os"
	"strings"
)

// Version is set at build time via -ldflags.
// Falls back to reading the VERSION file or "dev".
var Version = readVersion()

// ProductName is the display name.
const ProductName = "HyperHarness"

// GitCommit is set at build time.
var GitCommit = "unknown"

// BuildTime is set at build time.
var BuildTime = "unknown"

func readVersion() string {
	// Try to read from VERSION file in cwd or parent dirs
	for _, path := range []string{"VERSION", "../VERSION"} {
		if data, err := os.ReadFile(path); err == nil {
			v := strings.TrimSpace(string(data))
			if v != "" {
				return v
			}
		}
	}
	return "dev"
}
