package processor

import (
	"hdl/internal/archive"
	"hdl/internal/downloader"
	"hdl/internal/splitter"
)

type Processor struct {
	URL            string
	DownloadClient downloader.Downloader
	Archiver       archive.Archiver
	Splitter       splitter.Splitter
	OutputBase     string
}
