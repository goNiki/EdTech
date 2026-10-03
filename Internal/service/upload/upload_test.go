package upload_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"edtech/internal/service/upload"
	errorsAPP "edtech/pkg/errors"
)

type mockStorage struct {
	savedFiles map[string][]byte
}

func newMockStorage() *mockStorage {
	return &mockStorage{savedFiles: make(map[string][]byte)}
}

func (m *mockStorage) Save(ctx context.Context, relativePath string, src io.Reader) error {
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	m.savedFiles[relativePath] = data
	return nil
}

func (m *mockStorage) Delete(ctx context.Context, relativePath string) error {
	delete(m.savedFiles, relativePath)
	return nil
}

func (m *mockStorage) Exists(ctx context.Context, relativePath string) (bool, error) {
	_, ok := m.savedFiles[relativePath]
	return ok, nil
}

func TestUploadFile_Success(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	// Valid PNG bytes
	pngHeader := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15c4\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n-\xb4\x00\x00\x00\x00IEND\xaeB`\x82")
	reader := bytes.NewReader(pngHeader)

	res, err := svc.UploadFile(context.Background(), reader, "avatar.png", int64(len(pngHeader)), "avatar")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res == nil {
		t.Fatal("expected result, got nil")
	}

	if !strings.HasPrefix(res.FileURL, "/static/uploads/avatar/") {
		t.Errorf("expected URL prefix /static/uploads/avatar/, got %s", res.FileURL)
	}

	if res.MimeType != "image/png" {
		t.Errorf("expected mime image/png, got %s", res.MimeType)
	}

	if res.SizeBytes != int64(len(pngHeader)) {
		t.Errorf("expected size %d, got %d", len(pngHeader), res.SizeBytes)
	}

	if len(mockStore.savedFiles) != 1 {
		t.Errorf("expected 1 saved file in storage, got %d", len(mockStore.savedFiles))
	}
}

func TestUploadFile_PDF_Success(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	// Valid PDF header
	pdfContent := []byte("%PDF-1.4\n1 0 obj\n<<>>\nendobj\ntrailer\n<<>>\n%%EOF")
	reader := bytes.NewReader(pdfContent)

	res, err := svc.UploadFile(context.Background(), reader, "solution.pdf", int64(len(pdfContent)), "homework")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.HasPrefix(res.FileURL, "/static/uploads/homework/") {
		t.Errorf("expected URL prefix /static/uploads/homework/, got %s", res.FileURL)
	}

	if res.MimeType != "application/pdf" {
		t.Errorf("expected mime application/pdf, got %s", res.MimeType)
	}
}

func TestUploadFile_BlockedExecutableExtension(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	blockedFiles := []string{"malware.exe", "script.sh", "backdoor.php", "payload.js", "virus.bat"}

	for _, fname := range blockedFiles {
		data := []byte("some content")
		reader := bytes.NewReader(data)

		_, err := svc.UploadFile(context.Background(), reader, fname, int64(len(data)), "general")
		if err == nil {
			t.Errorf("expected error for blocked file %s, got nil", fname)
		}
		if !errors.Is(err, errorsAPP.ErrInvalidFileType) {
			t.Errorf("expected ErrInvalidFileType for %s, got %v", fname, err)
		}
	}
}

func TestUploadFile_TooLarge(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	data := []byte("small content")
	reader := bytes.NewReader(data)
	oversized := upload.MaxFileSize + 1

	_, err := svc.UploadFile(context.Background(), reader, "large.png", oversized, "general")
	if err == nil {
		t.Fatal("expected error for oversized file, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrFileTooLarge) {
		t.Fatalf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestUploadFile_EmptyFile(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	reader := bytes.NewReader([]byte{})
	_, err := svc.UploadFile(context.Background(), reader, "empty.png", 0, "general")
	if err == nil {
		t.Fatal("expected error for empty file, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrEmptyFile) {
		t.Fatalf("expected ErrEmptyFile, got %v", err)
	}
}
