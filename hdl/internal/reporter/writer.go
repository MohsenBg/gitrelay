package reporter

import (
	"encoding/json"
	"fmt"
	"os"
)

// WriteManifest writes the given manifest data to a JSON file at `path`.
func WriteManifest(path string, m Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}
	return nil
}
