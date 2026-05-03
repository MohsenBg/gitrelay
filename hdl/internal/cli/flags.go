package cli

import (
	"flag"
	"fmt"
	"os"

	"hdl/internal/processor"
	"hdl/internal/ui"
)

type Defaults struct {
	DefaultURL       string
	DefaultLinkType  string
	DefaultSplitType string
	DefaultFormat    string
	DefaultChunkSize string
	DefaultOutputDir string
}

var defaultConfig = Defaults{
	DefaultURL:       "",
	DefaultLinkType:  "direct",
	DefaultSplitType: "simple",
	DefaultFormat:    "zip",
	DefaultChunkSize: "95MB",
	DefaultOutputDir: "./downloads",
}

func ParseProcessorFlags() *processor.NewProcessorConfig {
	url := flag.String("url", defaultConfig.DefaultURL, "Download URL (required)")
	linkType := flag.String("link", defaultConfig.DefaultLinkType, "Link type (direct, youtube)")
	splitType := flag.String("split", defaultConfig.DefaultSplitType, "Split type (simple, archiver)")
	format := flag.String("format", defaultConfig.DefaultFormat, "Archive format (zip, tar)")
	chunkSizeStr := flag.String("chunksize", defaultConfig.DefaultChunkSize, "Max chunk size (<=100MB)")
	out := flag.String("out", defaultConfig.DefaultOutputDir, "Output directory")

	flag.Usage = func() {
		ui.Info("Usage: hdl -url <download_url> [options]")
		ui.PrintDefaults()
		fmt.Println()
	}

	flag.Parse()

	if *url == "" {
		flag.Usage()
		ui.Error("url is required")
		os.Exit(1)
	}

	size, err := parseSize(*chunkSizeStr)
	if err != nil {
		flag.Usage()
		ui.Error("invalid chunk size: %v", err)
		os.Exit(1)
	}

	const maxSize = 100 * 1000 * 1000
	if size > maxSize {
		flag.Usage()
		ui.Error("chunk size cannot exceed 100MB")
		os.Exit(1)
	}

	return &processor.NewProcessorConfig{
		URL:           *url,
		LinkType:      parseLinkType(*linkType),
		SplitType:     parseSplitType(*splitType),
		ArchiveFormat: parseArchiveFormat(*format),
		MaxChunkSize:  uint64(size),
		OutputBase:    *out,
	}
}

func printDefaults() {
	fmt.Println("Usage:")
	flag.PrintDefaults()

	fmt.Println("\nCurrent defaults:")
	fmt.Printf("  link      = %q\n", defaultConfig.DefaultLinkType)
	fmt.Printf("  split     = %q\n", defaultConfig.DefaultSplitType)
	fmt.Printf("  format    = %q\n", defaultConfig.DefaultFormat)
	fmt.Printf("  chunksize = %q\n", defaultConfig.DefaultChunkSize)
	fmt.Printf("  out       = %q\n", defaultConfig.DefaultOutputDir)
}
