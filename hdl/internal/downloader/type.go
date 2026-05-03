package downloader

import (
	"context"
	"fmt"
)

// LinkType represents the type of input link.
type LinkType int

const (
	LinkDirect LinkType = iota
	LinkYouTube
)

func (t LinkType) String() string {
	switch t {
	case LinkDirect:
		return "direct"
	case LinkYouTube:
		return "youtube"
	default:
		return fmt.Sprintf("LinkType(%d)", t)
	}
}

// Downloader defines the behavior for fetching remote content.
type Downloader interface {
	// Download fetches the file and returns the path to the downloaded file.
	Download(ctx context.Context, url, output string) (string, error)

	DownloadLinkType() LinkType
}
