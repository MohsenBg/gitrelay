package reporter

import (
	"time"
)

// FilePart represents metadata about a single file chunk.
type FilePart struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	Hash  string `json:"hash"`
	Size  int64  `json:"size"`
}

// ArchiveInfo stores details about the generated archive (zip/tar).
type ArchiveInfo struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	Hash   string `json:"hash"`
	Size   int64  `json:"size"`
}

// Manifest is the root metadata structure for a processed file set.
type Manifest struct {
	Original struct {
		Name string `json:"name"`
		Hash string `json:"sha256"`
	} `json:"original"`

	Archive ArchiveInfo `json:"archive"`

	ChunkSize  int64      `json:"chunk_size"`
	TotalParts int        `json:"total_parts"`
	CreatedAt  time.Time  `json:"created_at"`
	OutputPath string     `json:"output_location"`
	Parts      []FilePart `json:"parts"`
}

