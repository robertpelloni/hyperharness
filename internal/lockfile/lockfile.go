// Package lockfile provides process lock file management.
// Ported from hypercode/go/internal/lockfile/lockfile.go
package lockfile

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// StartupProvenance records how a process was started.
type StartupProvenance struct {
	RequestedRuntime string `json:"requestedRuntime,omitempty"`
	ActiveRuntime    string `json:"activeRuntime,omitempty"`
	RequestedPort    int    `json:"requestedPort,omitempty"`
	ActivePort       int    `json:"activePort,omitempty"`
	PortDecision     string `json:"portDecision,omitempty"`
	PortReason       string `json:"portReason,omitempty"`
	LaunchMode       string `json:"launchMode,omitempty"`
	DashboardMode    string `json:"dashboardMode,omitempty"`
	InstallDecision  string `json:"installDecision,omitempty"`
	InstallReason    string `json:"installReason,omitempty"`
	BuildDecision    string `json:"buildDecision,omitempty"`
	BuildReason      string `json:"buildReason,omitempty"`
	UpdatedAt        string `json:"updatedAt,omitempty"`
}

// Record holds the lock file data.
type Record struct {
	Host      string              `json:"host"`
	Port      int                 `json:"port"`
	Version   string              `json:"version"`
	StartedAt string              `json:"startedAt"`
	Startup   *StartupProvenance  `json:"startup,omitempty"`
}

// Write creates or overwrites a lock file.
func Write(lockPath string, record Record) error {
	if err := os.MkdirAll(filepath.Dir(lockPath), 0755); err != nil {
		return err
	}
	payload, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(lockPath, payload, 0644)
}

// Read reads a lock file.
func Read(lockPath string) (Record, error) {
	payload, err := os.ReadFile(lockPath)
	if err != nil {
		return Record{}, err
	}
	var record Record
	if err := json.Unmarshal(payload, &record); err != nil {
		return Record{}, err
	}
	return record, nil
}
