package storage

import (
	"context"
	"io"
)

// Storage defines the interface for persisting files (local disk, S3, MinIO)
type Storage interface {
	Save(ctx context.Context, relativePath string, src io.Reader) error
	Delete(ctx context.Context, relativePath string) error
	Exists(ctx context.Context, relativePath string) (bool, error)
}
