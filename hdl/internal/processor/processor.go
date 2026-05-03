package processor

import (
	"context"
	"fmt"
	"hdl/internal/archive"
	"hdl/internal/downloader"
	"hdl/internal/reporter"
	"hdl/internal/splitter"
	"hdl/internal/ui"
	"os"
	"path/filepath"
	"time"
)

type NewProcessorConfig struct {
	URL           string
	LinkType      downloader.LinkType
	SplitType     splitter.SplitType
	ArchiveFormat archive.ArchiveFormat
	MaxChunkSize  uint64
	OutputBase    string
}

func NewProcessor(config NewProcessorConfig) (*Processor, error) {
	if !isValidURL(config.URL) {
		return nil, fmt.Errorf("invalid url: %s", config.URL)
	}

	ui.Info("creating downloader (type=%v)", config.LinkType)
	dl, err := downloader.CreateDownloader(config.LinkType)
	if err != nil {
		return nil, fmt.Errorf("create downloader: %w", err)
	}

	ui.Info("creating splitter (type=%v, maxChunkSize=%d)", config.SplitType, config.MaxChunkSize)
	sp, err := splitter.CreateSplitter(config.SplitType, config.MaxChunkSize)
	if err != nil {
		return nil, fmt.Errorf("create splitter: %w", err)
	}

	ui.Info("creating archiver (format=%v)", config.ArchiveFormat)
	ar, err := archive.CreateArchiver(config.ArchiveFormat)
	if err != nil {
		return nil, fmt.Errorf("create archiver: %w", err)
	}

	ui.Success("processor created")

	return &Processor{
		URL:            config.URL,
		DownloadClient: dl,
		Archiver:       ar,
		Splitter:       sp,
		OutputBase:     config.OutputBase,
	}, nil
}

func (p *Processor) Run(ctx context.Context) error {
	ui.Info("starting processor for %s", p.URL)

	if err := ctx.Err(); err != nil {
		ui.Warning("context cancelled: %v", err)
		return err
	}

	return p.run(ctx)
}

func (p *Processor) run(ctx context.Context) error {
	tmpDir, cleanup, err := p.createTempDir()
	if err != nil {
		return err
	}
	defer cleanup()

	downloadPath, err := p.download(ctx, tmpDir)
	if err != nil {
		return err
	}

	origHash, err := p.hashOriginal(downloadPath)
	if err != nil {
		return err
	}

	if p.isAlreadyProcessed(origHash) {
		ui.Success("already processed, using cached output: %s", origHash)
		return nil
	}

	partsDir, err := p.prepareOutputDir(origHash)
	if err != nil {
		return err
	}

	archiveFile, archiveInfo, archiveHash, err := p.compress(downloadPath, tmpDir)
	if err != nil {
		return err
	}

	parts, err := p.splitArchive(archiveFile, partsDir)
	if err != nil {
		return err
	}

	return p.writeManifest(
		downloadPath,
		partsDir,
		archiveFile,
		archiveInfo,
		archiveHash,
		origHash,
		parts,
	)
}

func (p *Processor) createTempDir() (string, func(), error) {
	ui.Info("creating temporary working directory")

	tmpDir, err := os.MkdirTemp("", "download_*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp directory: %w", err)
	}

	ui.Success("temp directory created: %s", tmpDir)

	return tmpDir, func() { _ = os.RemoveAll(tmpDir) }, nil
}

func (p *Processor) download(ctx context.Context, tmpDir string) (string, error) {
	ui.Info("downloading file: %s", p.URL)

	path, err := p.DownloadClient.Download(ctx, p.URL, tmpDir)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}

	ui.Success("download completed: %s", path)
	return path, nil
}

func (p *Processor) prepareOutputDir(hash string) (string, error) {
	dir := filepath.Join(p.OutputBase, hash)

	ui.Info("creating output directory: %s", dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}

	ui.Success("output directory ready")
	return dir, nil
}

func (p *Processor) compress(src, tmpDir string) (string, os.FileInfo, string, error) {
	ui.Info("compressing file")

	archiveFile, err := p.Archiver.Compress(src, tmpDir)
	if err != nil {
		return "", nil, "", fmt.Errorf("compression failed: %w", err)
	}

	info, err := os.Stat(archiveFile)
	if err != nil {
		return "", nil, "", fmt.Errorf("stat archive file: %w", err)
	}

	ui.Info("calculating SHA256 of archive")
	hash, err := CalcFileSHA256(archiveFile)
	if err != nil {
		return "", nil, "", fmt.Errorf("calculate archive hash: %w", err)
	}

	ui.Success("archive created: %s", archiveFile)
	ui.Success("archive hash: %s", hash)

	return archiveFile, info, hash, nil
}

func (p *Processor) splitArchive(archiveFile, outDir string) ([]string, error) {
	ui.Info("splitting archive")

	parts, err := p.Splitter.SplitFile(archiveFile, outDir, 0)
	if err != nil {
		return nil, fmt.Errorf("file splitting failed: %w", err)
	}

	ui.Success("archive split into %d parts", len(parts))
	return parts, nil
}

func (p *Processor) writeManifest(
	downloadPath string,
	partsDir string,
	archiveFile string,
	archiveInfo os.FileInfo,
	archiveHash string,
	origHash string,
	parts []string,
) error {
	ui.Info("generating manifest")

	manifest := reporter.Manifest{
		CreatedAt:  time.Now(),
		OutputPath: partsDir,
		ChunkSize:  int64(p.Splitter.MaxChunkSize()),
		TotalParts: len(parts),
	}

	manifest.Original.Name = filepath.Base(downloadPath)
	manifest.Original.Hash = origHash

	manifest.Archive = reporter.ArchiveInfo{
		Name:   filepath.Base(archiveFile),
		Format: string(p.Archiver.Format()),
		Hash:   archiveHash,
		Size:   archiveInfo.Size(),
	}

	for i, part := range parts {
		info, err := os.Stat(part)
		if err != nil {
			ui.Warning("stat failed for %s: %v", part, err)
			continue
		}

		hash, err := CalcFileSHA256(part)
		if err != nil {
			ui.Warning("hash failed for %s: %v", part, err)
			continue
		}

		manifest.Parts = append(manifest.Parts, reporter.FilePart{
			Index: i + 1,
			Name:  filepath.Base(part),
			Hash:  hash,
			Size:  info.Size(),
		})
	}

	path := filepath.Join(partsDir, "manifest.json")
	ui.Info("writing manifest: %s", path)

	if err := reporter.WriteManifest(path, manifest); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}

	ui.Success("manifest written")
	return nil
}

func (p *Processor) hashOriginal(path string) (string, error) {
	ui.Info("calculating SHA256 of original file")

	hash, err := CalcFileSHA256(path)
	if err != nil {
		return "", fmt.Errorf("calculate original hash: %w", err)
	}

	ui.Success("original file hash: %s", hash)
	return hash, nil
}

func (p *Processor) isAlreadyProcessed(hash string) bool {
	hashDir := filepath.Join(p.OutputBase, hash)

	manifestPath := filepath.Join(hashDir, "manifest.json")
	manifest, err := reporter.ReadManifest(manifestPath)
	if err != nil {
		return false
	}

	files, err := os.ReadDir(hashDir)
	if err != nil {
		return false
	}

	foundParts := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		matched, _ := filepath.Match("*.[0-9][0-9][0-9]", name)
		if matched {
			foundParts++
		}
	}

	return foundParts == manifest.TotalParts
}
