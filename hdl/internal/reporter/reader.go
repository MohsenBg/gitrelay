package reporter

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadManifest reads and parses a manifest JSON file from `path` into a Manifest struct.
func ReadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}

	return &m, nil
}
