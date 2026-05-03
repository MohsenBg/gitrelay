package main

import (
	"context"
	"os"

	"hdl/internal/cli"
	"hdl/internal/processor"
	"hdl/internal/ui"
)

func main() {
	ui.Init()

	cfg := cli.ParseProcessorFlags()

	p, err := processor.NewProcessor(*cfg)
	if err != nil {
		ui.Error("failed to initialize processor: %v", err)
		os.Exit(1)
	}

	if err := p.Run(context.Background()); err != nil {
		ui.Error("processing failed: %v", err)
		os.Exit(1)
	}

	ui.Success("processing completed successfully")
}
