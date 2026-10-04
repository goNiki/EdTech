package upload_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/service/upload"
	errorsAPP "edtech/pkg/errors"
)

type mockStorage struct {
	mu         sync.Mutex
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
	m.mu.Lock()
	defer m.mu.Unlock()
	m.savedFiles[relativePath] = data
	return nil
}

func (m *mockStorage) Delete(ctx context.Context, relativePath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.savedFiles, relativePath)
	return nil
}

func (m *mockStorage) Exists(ctx context.Context, relativePath string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.savedFiles[relativePath]
	return ok, nil
}

func (m *mockStorage) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.savedFiles)
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

	if mockStore.Count() != 1 {
		t.Errorf("expected 1 saved file in storage, got %d", mockStore.Count())
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

func TestUploadImagesBatch_Success(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	pngBytes := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15c4\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n-\xb4\x00\x00\x00\x00IEND\xaeB`\x82")
	jpegBytes := []byte("\xff\xd8\xff\xe0\x00\x10JFIF\x00\x01\x01\x00\x00\x01\x00\x01\x00\x00\xff\xdb\x00C\x00\x08\x06\x06\x07\x06\x05\x08\x07\x07\x07\t\t\x08\n\x0c\x14\r\x0c\x0b\x0b\x0c\x19\x12\x13\x0f\x14\x1d\x1a\x1f\x1e\x1d\x1a\x1c\x1c $.' \",#\x1c\x1c(7),01444\x1f'9=82<.342\xff\xd9")

	files := []domain.BatchFileItem{
		{
			Reader:   bytes.NewReader(pngBytes),
			Filename: "diag1.png",
			Size:     int64(len(pngBytes)),
		},
		{
			Reader:   bytes.NewReader(jpegBytes),
			Filename: "chart2.jpg",
			Size:     int64(len(jpegBytes)),
		},
	}

	results, err := svc.UploadImagesBatch(context.Background(), files, "lesson_media")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].OriginalName != "diag1.png" || !strings.HasPrefix(results[0].FileURL, "/static/uploads/lesson_media/") {
		t.Errorf("invalid result 0: %+v", results[0])
	}
	if results[1].OriginalName != "chart2.jpg" || !strings.HasPrefix(results[1].FileURL, "/static/uploads/lesson_media/") {
		t.Errorf("invalid result 1: %+v", results[1])
	}

	if mockStore.Count() != 2 {
		t.Errorf("expected 2 saved files in mockStore, got %d", mockStore.Count())
	}
}

func TestUploadImagesBatch_TooManyFiles(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	files := make([]domain.BatchFileItem, 51)
	for i := 0; i < 51; i++ {
		files[i] = domain.BatchFileItem{
			Reader:   bytes.NewReader([]byte("fake")),
			Filename: fmt.Sprintf("img%d.png", i),
			Size:     4,
		}
	}

	_, err := svc.UploadImagesBatch(context.Background(), files, "lesson_media")
	if err == nil {
		t.Fatal("expected error for 51 files, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrBatchTooManyFiles) {
		t.Fatalf("expected ErrBatchTooManyFiles, got %v", err)
	}
}

func TestUploadImagesBatch_TotalSizeTooLarge(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	// 2 files with total size > 50 MB
	files := []domain.BatchFileItem{
		{
			Reader:   bytes.NewReader([]byte("1")),
			Filename: "img1.png",
			Size:     25 * 1024 * 1024,
		},
		{
			Reader:   bytes.NewReader([]byte("2")),
			Filename: "img2.png",
			Size:     25*1024*1024 + 1,
		},
	}

	_, err := svc.UploadImagesBatch(context.Background(), files, "lesson_media")
	if err == nil {
		t.Fatal("expected error for total size > 50MB, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrFileTooLarge) {
		t.Fatalf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestUploadImagesBatch_NonImageRejected(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	pdfContent := []byte("%PDF-1.4\n1 0 obj\n<<>>\nendobj\ntrailer\n<<>>\n%%EOF")
	files := []domain.BatchFileItem{
		{
			Reader:   bytes.NewReader(pdfContent),
			Filename: "document.pdf",
			Size:     int64(len(pdfContent)),
		},
	}

	_, err := svc.UploadImagesBatch(context.Background(), files, "lesson_media")
	if err == nil {
		t.Fatal("expected error for PDF in images batch, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrInvalidFileType) {
		t.Fatalf("expected ErrInvalidFileType, got %v", err)
	}
}

func TestUploadImagesBatch_EmptyBatch(t *testing.T) {
	mockStore := newMockStorage()
	svc := upload.NewUploadService(mockStore)

	_, err := svc.UploadImagesBatch(context.Background(), nil, "lesson_media")
	if err == nil {
		t.Fatal("expected error for empty batch, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrEmptyFile) {
		t.Fatalf("expected ErrEmptyFile, got %v", err)
	}
}
