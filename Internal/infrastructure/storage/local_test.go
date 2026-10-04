package storage_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"edtech/internal/infrastructure/storage"
	errorsAPP "edtech/pkg/errors"
)

func TestLocalStorage_ValidPaths(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "storage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	s, err := storage.NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("failed to init local storage: %v", err)
	}

	ctx := context.Background()
	relPath := filepath.Join("avatar", "2026", "04", "test.txt")
	content := []byte("hello world")

	// 1. Exists before save -> false
	exists, err := s.Exists(ctx, relPath)
	if err != nil {
		t.Fatalf("unexpected error checking existence: %v", err)
	}
	if exists {
		t.Errorf("expected exists=false before saving")
	}

	// 2. Save
	if err := s.Save(ctx, relPath, bytes.NewReader(content)); err != nil {
		t.Fatalf("failed to save file: %v", err)
	}

	// 3. Exists after save -> true
	exists, err = s.Exists(ctx, relPath)
	if err != nil {
		t.Fatalf("unexpected error checking existence: %v", err)
	}
	if !exists {
		t.Errorf("expected exists=true after saving")
	}

	// Verify content written on disk
	diskPath := filepath.Join(tempDir, relPath)
	readBytes, err := os.ReadFile(diskPath)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(readBytes) != string(content) {
		t.Errorf("expected content %q, got %q", string(content), string(readBytes))
	}

	// 4. Delete
	if err := s.Delete(ctx, relPath); err != nil {
		t.Fatalf("failed to delete file: %v", err)
	}

	// 5. Exists after delete -> false
	exists, err = s.Exists(ctx, relPath)
	if err != nil {
		t.Fatalf("unexpected error checking existence: %v", err)
	}
	if exists {
		t.Errorf("expected exists=false after deletion")
	}
}

func TestLocalStorage_PathTraversalAttacks(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "storage_traversal_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	s, err := storage.NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("failed to init local storage: %v", err)
	}

	ctx := context.Background()

	maliciousPayloads := []string{
		"../../etc/passwd",
		"..\\..\\Windows\\System32\\cmd.exe",
		"../../../secret.txt",
		"avatar/../../../etc/shadow",
		"/etc/passwd",
		"\\Windows\\System32\\config\\sam",
		"C:\\Windows\\System32\\calc.exe",
		"C:/Windows/System32/calc.exe",
		"D:\\data.txt",
		"..",
		"../",
		"..\\",
		".",
		"",
		"   ",
	}

	for _, payload := range maliciousPayloads {
		t.Run("Payload_"+payload, func(t *testing.T) {
			// Test Save
			err := s.Save(ctx, payload, bytes.NewReader([]byte("evil")))
			if err == nil {
				t.Fatalf("expected error on Save with payload %q, got nil", payload)
			}
			if !errors.Is(err, errorsAPP.ErrPathTraversal) {
				t.Errorf("expected ErrPathTraversal on Save for %q, got: %v", payload, err)
			}

			// Test Exists
			exists, err := s.Exists(ctx, payload)
			if err == nil {
				t.Fatalf("expected error on Exists with payload %q, got nil", payload)
			}
			if exists {
				t.Errorf("expected exists=false for %q", payload)
			}
			if !errors.Is(err, errorsAPP.ErrPathTraversal) {
				t.Errorf("expected ErrPathTraversal on Exists for %q, got: %v", payload, err)
			}

			// Test Delete
			err = s.Delete(ctx, payload)
			if err == nil {
				t.Fatalf("expected error on Delete with payload %q, got nil", payload)
			}
			if !errors.Is(err, errorsAPP.ErrPathTraversal) {
				t.Errorf("expected ErrPathTraversal on Delete for %q, got: %v", payload, err)
			}
		})
	}
}
