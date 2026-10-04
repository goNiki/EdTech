package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	errorsAPP "edtech/pkg/errors"
)

type LocalStorage struct {
	baseDir string
}

func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("create base dir %s: %w", baseDir, err)
	}
	return &LocalStorage{baseDir: baseDir}, nil
}

func (s *LocalStorage) resolveSafePath(relativePath string) (string, error) {
	trimmed := strings.TrimSpace(relativePath)
	if trimmed == "" {
		slog.Warn("storage: empty relative path")
		return "", fmt.Errorf("%w: path cannot be empty", errorsAPP.ErrPathTraversal)
	}

	cleanBase, err := filepath.Abs(s.baseDir)
	if err != nil {
		return "", fmt.Errorf("resolve base dir: %w", err)
	}

	// Reject absolute paths, volume drive specifications, or rooted slashes upfront
	if filepath.IsAbs(trimmed) || filepath.VolumeName(trimmed) != "" || strings.HasPrefix(trimmed, "/") || strings.HasPrefix(trimmed, "\\") {
		slog.Warn("storage: path traversal attempt rejected (rooted/volume path)", "path", relativePath)
		return "", fmt.Errorf("%w: rooted path or drive letter not allowed: %s", errorsAPP.ErrPathTraversal, relativePath)
	}

	targetPath := filepath.Clean(filepath.Join(cleanBase, trimmed))
	rel, err := filepath.Rel(cleanBase, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") || rel == "." {
		slog.Warn("storage: path traversal attempt rejected", "path", relativePath, "targetPath", targetPath, "rel", rel)
		return "", fmt.Errorf("%w: path traversal detected for '%s'", errorsAPP.ErrPathTraversal, relativePath)
	}

	return targetPath, nil
}

func (s *LocalStorage) Save(ctx context.Context, relativePath string, src io.Reader) error {
	const op = "infrastructure.storage.local.Save"

	fullPath, err := s.resolveSafePath(relativePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%s: create directory %s: %w", op, dir, err)
	}

	dst, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("%s: create file %s: %w", op, fullPath, err)
	}
	defer func() {
		_ = dst.Close()
	}()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("%s: copy content: %w", op, err)
	}

	return nil
}

func (s *LocalStorage) Delete(ctx context.Context, relativePath string) error {
	const op = "infrastructure.storage.local.Delete"

	fullPath, err := s.resolveSafePath(relativePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%s: delete file %s: %w", op, fullPath, err)
	}
	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, relativePath string) (bool, error) {
	const op = "infrastructure.storage.local.Exists"

	fullPath, err := s.resolveSafePath(relativePath)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	_, err = os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
