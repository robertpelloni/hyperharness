package sessionimport

import "time"

func BuildManifest(candidates []ValidationResult) Manifest {
	return Manifest{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Count:       len(candidates),
		Candidates:  candidates,
	}
}
