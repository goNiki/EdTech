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

	"edtech/internal/domain"
	"edtech/internal/infrastructure/storage"
	services "edtech/internal/service"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
)

const (
	MaxFileSize int64 = 25 * 1024 * 1024 // 25 MB
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

var validCategories = map[string]bool{
	"avatar":       true,
	"course_cover": true,
	"homework":     true,
	"general":      true,
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

	if size > MaxFileSize {
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

	category = strings.ToLower(strings.TrimSpace(category))
	if !validCategories[category] {
		category = "general"
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
