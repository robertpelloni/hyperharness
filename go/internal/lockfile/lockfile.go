package lockfile

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Record struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Version   string `json:"version"`
	StartedAt string `json:"startedAt"`
}

func Write(lockPath string, record Record) error {
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(lockPath, payload, 0o644)
}

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
