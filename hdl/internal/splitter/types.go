package splitter

import "fmt"

type Splitter interface {
	SplitFile(sourcePath, baseTargetPath string, chunkSize uint64) ([]string, error)
	MergeFiles(basePath, outputPath string) error
	MaxChunkSize() uint64
}

type SimpleSplitter struct {
	chunkSize uint64
}

// SplitType represents supported splitting strategies.
type SplitType int

const (
	SplitSimple SplitType = iota
	SplitArchiver
)

func (s SplitType) String() string {
	switch s {
	case SplitSimple:
		return "simple"
	case SplitArchiver:
		return "archiver"
	default:
		return fmt.Sprintf("SplitType(%d)", s)
	}
}
