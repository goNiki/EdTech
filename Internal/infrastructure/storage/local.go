package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func (s *LocalStorage) Save(ctx context.Context, relativePath string, src io.Reader) error {
	const op = "infrastructure.storage.local.Save"

	fullPath := filepath.Join(s.baseDir, filepath.Clean(relativePath))
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%s: create directory %s: %w", op, dir, err)
	}

	dst, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("%s: create file %s: %w", op, fullPath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("%s: copy content: %w", op, err)
	}

	return nil
}

func (s *LocalStorage) Delete(ctx context.Context, relativePath string) error {
	const op = "infrastructure.storage.local.Delete"

	fullPath := filepath.Join(s.baseDir, filepath.Clean(relativePath))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%s: delete file %s: %w", op, fullPath, err)
	}
	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, relativePath string) (bool, error) {
	fullPath := filepath.Join(s.baseDir, filepath.Clean(relativePath))
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
