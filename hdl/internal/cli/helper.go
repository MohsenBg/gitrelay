package cli

import (
	"fmt"
	"hdl/internal/archive"
	"hdl/internal/downloader"
	"hdl/internal/splitter"
	"hdl/internal/ui"
	"os"
	"strings"
)

func parseLinkType(v string) downloader.LinkType {
	switch v {
	case "direct":
		return downloader.LinkDirect
	case "youtube":
		return downloader.LinkYouTube
	default:
		ui.Error("invalid link type: %s (allowed: direct, youtube)", v)
		os.Exit(1)
	}
	return downloader.LinkDirect
}

func parseSplitType(v string) splitter.SplitType {
	switch v {
	case "simple":
		return splitter.SplitSimple
	case "archiver":
		return splitter.SplitArchiver
	default:
		ui.Error("invalid split type: %s (allowed: simple, archiver)", v)
		os.Exit(1)
	}
	return splitter.SplitSimple
}

func parseArchiveFormat(v string) archive.ArchiveFormat {
	switch v {
	case "zip":
		return archive.ArchiveZIP
	case "tar":
		return archive.ArchiveTAR
	default:
		ui.Error("invalid archive format: %s (allowed: zip, tar)", v)
		os.Exit(1)
	}
	return archive.ArchiveZIP
}

func parseSize(s string) (uint64, error) {
	var num float64
	var unit string

	_, err := fmt.Sscanf(strings.ToUpper(strings.TrimSpace(s)), "%f%s", &num, &unit)
	if err != nil {
		ui.Error("invalid size format")
		os.Exit(1)
	}

	switch unit {
	case "B", "":
		return uint64(num), nil
	case "KB":
		return uint64(num * 1000), nil
	case "MB":
		return uint64(num * 1000 * 1000), nil
	case "GB":
		return uint64(num * 1000 * 1000 * 1000), nil
	default:
		return 0, fmt.Errorf("unknown unit %s", unit)
	}
}
