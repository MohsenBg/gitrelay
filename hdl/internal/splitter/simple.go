package splitter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func NewSimpleSplitter(maxChunkSize uint64) (Splitter, error) {
	if maxChunkSize == 0 {
		return nil, fmt.Errorf("maxChunkSize must be > 0")
	}

	return &SimpleSplitter{chunkSize: maxChunkSize}, nil
}

func (s *SimpleSplitter) SplitFile(sourcePath, targetDir string, chunkSize uint64) ([]string, error) {
	if chunkSize == 0 {
		chunkSize = s.chunkSize
	}
	if chunkSize > s.chunkSize {
		return nil, fmt.Errorf("chunkSize %d exceeds max allowed %d", chunkSize, s.chunkSize)
	}

	in, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("open source file: %w", err)
	}
	defer in.Close()

	base := filepath.Join(targetDir, filepath.Base(sourcePath))
	return s.split(in, base, chunkSize)
}

func (s *SimpleSplitter) MaxChunkSize() uint64 {
	return s.chunkSize
}

func (s *SimpleSplitter) split(in *os.File, basePath string, chunkSize uint64) ([]string, error) {
	var paths []string

	for partNum := 1; ; partNum++ {
		partPath := fmt.Sprintf("%s.%03d", basePath, partNum)

		out, err := os.Create(partPath)
		if err != nil {
			return nil, fmt.Errorf("create part %q: %w", partPath, err)
		}

		n, err := io.CopyN(out, in, int64(chunkSize))
		out.Close()

		if err == io.EOF {
			if n == 0 {
				os.Remove(partPath)
				break
			}
			paths = append(paths, partPath)
			break
		}

		if err != nil {
			return nil, fmt.Errorf("write part %q: %w", partPath, err)
		}

		paths = append(paths, partPath)
	}

	return paths, nil
}

// MergeFiles merges parts like basePath.001, .002, ... into outputPath.
func (s *SimpleSplitter) MergeFiles(basePath, outputPath string) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer out.Close()

	var found bool

	for partNum := 1; ; partNum++ {
		partPath := fmt.Sprintf("%s.%03d", basePath, partNum)

		in, err := os.Open(partPath)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return fmt.Errorf("open part %q: %w", partPath, err)
		}

		if _, err := io.Copy(out, in); err != nil {
			in.Close()
			return fmt.Errorf("merge part %q: %w", partPath, err)
		}

		in.Close()
		found = true
	}

	if !found {
		return fmt.Errorf("no parts found for base %s", basePath)
	}

	return nil
}
