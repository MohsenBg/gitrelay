package splitter

import "fmt"

// CreateSplitter initializes a splitter based on the given SplitType and max chunk size.
func CreateSplitter(splitType SplitType, maxChunkSize uint64) (Splitter, error) {
	switch splitType {
	case SplitSimple:
		return NewSimpleSplitter(maxChunkSize)
	case SplitArchiver:
		return nil, fmt.Errorf("archiver splitter not supported yet")
	default:
		return nil, fmt.Errorf("unsupported split type: %s", splitType.String())
	}
}
