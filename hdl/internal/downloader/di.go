package downloader

import "fmt"

// createDownloader initializes a downloader based on the given LinkType.
func CreateDownloader(linkType LinkType) (Downloader, error) {
	switch linkType {
	case LinkDirect:
		return NewDirectDownloader(), nil
	case LinkYouTube:
		return nil, fmt.Errorf("youtube downloader not supported yet")
	default:
		return nil, fmt.Errorf("unsupported link type: %s", linkType.String())
	}
}
