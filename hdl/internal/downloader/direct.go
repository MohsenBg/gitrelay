package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type DirectDownloader struct{}

// NewDirectDownloader creates a new instance of DirectDownloader.
func NewDirectDownloader() Downloader {
	return &DirectDownloader{}
}

func (d *DirectDownloader) DownloadLinkType() LinkType {
	return LinkDirect
}

// Download fetches the file from the URL provided in the config and saves it to disk.
func (d *DirectDownloader) Download(ctx context.Context, url, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	fileName := filepath.Base(url)
	if fileName == "." || fileName == "/" {
		fileName = "downloaded_file"
	}

	targetPath := filepath.Join(outputDir, fileName)

	out, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	return targetPath, nil
}
