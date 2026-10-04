package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/storage"
	services "edtech/internal/service"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
)

const (
	MaxFileSize             int64 = 25 * 1024 * 1024 // 25 MB
	MaxPresentationFileSize int64 = 50 * 1024 * 1024 // 50 MB
	MaxBatchFileCount       int   = 50               // 50 files max per batch
	MaxBatchTotalSize       int64 = 50 * 1024 * 1024 // 50 MB max total size per batch
)

var _ services.UploadServices = (*uploadService)(nil)

var dangerousExtensions = map[string]bool{
	".exe":  true,
	".sh":   true,
	".php":  true,
	".js":   true,
	".bat":  true,
	".cmd":  true,
	".bin":  true,
	".msi":  true,
	".dll":  true,
	".com":  true,
	".vbs":  true,
	".ps1":  true,
	".py":   true,
	".html": true,
	".htm":  true,
	".jar":  true,
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg":                   true,
	"image/png":                    true,
	"image/webp":                   true,
	"image/gif":                    true,
	"application/pdf":              true,
	"application/zip":              true,
	"application/x-zip-compressed": true,
}

var allowedBatchImageMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

var allowedBatchImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

var validCategories = map[string]bool{
	"avatar":       true,
	"course_cover": true,
	"homework":     true,
	"general":      true,
	"lesson_media": true,
	"presentation": true,
}

type uploadService struct {
	storage storage.Storage
}

func NewUploadService(strg storage.Storage) *uploadService {
	return &uploadService{
		storage: strg,
	}
}

func (s *uploadService) UploadFile(ctx context.Context, file io.Reader, filename string, size int64, category string) (*domain.FileUploadResult, error) {
	const op = "service.upload.UploadFile"

	if size <= 0 {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrEmptyFile)
	}

	category = strings.ToLower(strings.TrimSpace(category))
	if !validCategories[category] {
		category = "general"
	}

	maxAllowedSize := MaxFileSize
	if category == "presentation" {
		maxAllowedSize = MaxPresentationFileSize
	}

	if size > maxAllowedSize {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrFileTooLarge)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if dangerousExtensions[ext] {
		slog.Warn("attempt to upload dangerous file extension",
			slog.String("filename", filename),
			slog.String("extension", ext),
		)
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidFileType)
	}

	if category == "presentation" && ext != ".pdf" {
		slog.Warn("attempt to upload non-pdf extension for presentation category",
			slog.String("filename", filename),
			slog.String("extension", ext),
		)
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidFileType)
	}

	// Read first 512 bytes for MIME sniffing
	buf := make([]byte, 512)
	n, err := io.ReadFull(file, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("%s: read header bytes: %w", op, err)
	}
	if n == 0 {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrEmptyFile)
	}

	detectedMime := http.DetectContentType(buf[:n])

	// Special case: zip files can be detected as application/octet-stream or application/zip
	if (detectedMime == "application/octet-stream" || detectedMime == "application/zip") && (ext == ".zip") {
		detectedMime = "application/zip"
	}

	// Normalize mime type (remove charset if present)
	mimeClean := strings.Split(detectedMime, ";")[0]
	if !allowedMimeTypes[mimeClean] {
		slog.Warn("attempt to upload unsupported MIME type",
			slog.String("filename", filename),
			slog.String("detected_mime", detectedMime),
		)
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidFileType)
	}

	if category == "presentation" && mimeClean != "application/pdf" {
		slog.Warn("attempt to upload non-pdf MIME type for presentation category",
			slog.String("filename", filename),
			slog.String("detected_mime", detectedMime),
		)
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidFileType)
	}

	uuidStr, err := utils.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("%s: generate uuid: %w", op, err)
	}

	now := time.Now()
	yearMonth := now.Format("2006/01")
	newFilename := fmt.Sprintf("%s%s", uuidStr, ext)
	relativePath := filepath.Join(category, yearMonth, newFilename)
	// Convert Windows path separators to URL path slashes
	urlSubPath := fmt.Sprintf("%s/%s/%s", category, yearMonth, newFilename)

	fullReader := io.MultiReader(bytes.NewReader(buf[:n]), file)

	if err := s.storage.Save(ctx, relativePath, fullReader); err != nil {
		return nil, fmt.Errorf("%s: save file: %w", op, err)
	}

	slog.Info("file successfully uploaded",
		slog.String("original_name", filename),
		slog.String("saved_path", relativePath),
		slog.Int64("size_bytes", size),
		slog.String("mime_type", mimeClean),
		slog.String("category", category),
	)

	fileURL := fmt.Sprintf("/static/uploads/%s", urlSubPath)

	return &domain.FileUploadResult{
		FileURL:   fileURL,
		FileName:  filename,
		SizeBytes: size,
		MimeType:  mimeClean,
	}, nil
}

func (s *uploadService) UploadImagesBatch(ctx context.Context, files []domain.BatchFileItem, category string) ([]domain.BatchUploadResultItem, error) {
	const op = "service.upload.UploadImagesBatch"

	if len(files) == 0 {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrEmptyFile)
	}

	if len(files) > MaxBatchFileCount {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrBatchTooManyFiles)
	}

	var totalSize int64
	for _, f := range files {
		if f.Size <= 0 {
			return nil, fmt.Errorf("%s: file %q is empty: %w", op, f.Filename, errorsAPP.ErrEmptyFile)
		}
		if f.Size > MaxFileSize {
			return nil, fmt.Errorf("%s: file %q exceeds 25MB limit: %w", op, f.Filename, errorsAPP.ErrFileTooLarge)
		}
		totalSize += f.Size
	}

	if totalSize > MaxBatchTotalSize {
		return nil, fmt.Errorf("%s: total batch size %d exceeds 50MB limit: %w", op, totalSize, errorsAPP.ErrFileTooLarge)
	}

	category = strings.ToLower(strings.TrimSpace(category))
	if !validCategories[category] {
		category = "lesson_media"
	}

	now := time.Now()
	yearMonth := now.Format("2006/01")

	g, gCtx := errgroup.WithContext(ctx)
	results := make([]domain.BatchUploadResultItem, len(files))
	sem := make(chan struct{}, 10)

	for i, item := range files {
		i := i
		fItem := item
		g.Go(func() error {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-gCtx.Done():
				return gCtx.Err()
			}

			ext := strings.ToLower(filepath.Ext(fItem.Filename))
			if dangerousExtensions[ext] || !allowedBatchImageExtensions[ext] {
				slog.Warn("attempt to upload invalid image extension in batch",
					slog.String("filename", fItem.Filename),
					slog.String("extension", ext),
				)
				return fmt.Errorf("%s: file %q: %w", op, fItem.Filename, errorsAPP.ErrInvalidFileType)
			}

			buf := make([]byte, 512)
			n, err := io.ReadFull(fItem.Reader, buf)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return fmt.Errorf("%s: read header for %q: %w", op, fItem.Filename, err)
			}
			if n == 0 {
				return fmt.Errorf("%s: empty content for %q: %w", op, fItem.Filename, errorsAPP.ErrEmptyFile)
			}

			detectedMime := http.DetectContentType(buf[:n])
			mimeClean := strings.Split(detectedMime, ";")[0]
			if !allowedBatchImageMimeTypes[mimeClean] {
				slog.Warn("attempt to upload unsupported image MIME type in batch",
					slog.String("filename", fItem.Filename),
					slog.String("detected_mime", detectedMime),
				)
				return fmt.Errorf("%s: file %q has unsupported mime %q: %w", op, fItem.Filename, detectedMime, errorsAPP.ErrInvalidFileType)
			}

			uuidStr, err := utils.GenerateUUID()
			if err != nil {
				return fmt.Errorf("%s: generate uuid: %w", op, err)
			}

			newFilename := fmt.Sprintf("%s%s", uuidStr, ext)
			relativePath := filepath.Join(category, yearMonth, newFilename)
			urlSubPath := fmt.Sprintf("%s/%s/%s", category, yearMonth, newFilename)

			fullReader := io.MultiReader(bytes.NewReader(buf[:n]), fItem.Reader)
			if err := s.storage.Save(gCtx, relativePath, fullReader); err != nil {
				return fmt.Errorf("%s: save %q: %w", op, fItem.Filename, err)
			}

			results[i] = domain.BatchUploadResultItem{
				OriginalName: fItem.Filename,
				FileURL:      fmt.Sprintf("/static/uploads/%s", urlSubPath),
				SizeBytes:    fItem.Size,
				MimeType:     mimeClean,
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	slog.Info("batch images successfully uploaded",
		slog.Int("count", len(results)),
		slog.String("category", category),
		slog.Int64("total_bytes", totalSize),
	)

	return results, nil
}

