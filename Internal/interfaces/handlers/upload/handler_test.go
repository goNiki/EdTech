package upload_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/dto"
	uploadHandler "edtech/internal/interfaces/handlers/upload"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"
)

type mockUploadService struct {
	service.UploadServices
	uploadedFile *domain.FileUploadResult
	batchResults []domain.BatchUploadResultItem
	err          error
}

func (m *mockUploadService) UploadFile(ctx context.Context, file io.Reader, filename string, size int64, category string) (*domain.FileUploadResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.uploadedFile, nil
}

func (m *mockUploadService) UploadImagesBatch(ctx context.Context, files []domain.BatchFileItem, category string) ([]domain.BatchUploadResultItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.batchResults, nil
}

type mockAuthMiddleware struct {
	auth.AuthMiddleware
	userID int64
}

func (m *mockAuthMiddleware) GetUserID(ctx context.Context) int64 {
	return m.userID
}

type zeroReader struct {
	remaining int64
}

func (z *zeroReader) Read(p []byte) (n int, err error) {
	if z.remaining <= 0 {
		return 0, io.EOF
	}
	toRead := int64(len(p))
	if toRead > z.remaining {
		toRead = z.remaining
	}
	for i := int64(0); i < toRead; i++ {
		p[i] = 'X'
	}
	z.remaining -= toRead
	return int(toRead), nil
}

func TestUploadFile_MaxBytesReader_Protection(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	authMw := &mockAuthMiddleware{userID: 1}
	uploadSvc := &mockUploadService{
		uploadedFile: &domain.FileUploadResult{
			FileURL:   "/uploads/avatar/file.png",
			FileName:  "file.png",
			SizeBytes: 1024,
			MimeType:  "image/png",
		},
	}

	h := uploadHandler.NewUploadHandler(uploadSvc, authMw, logger)

	t.Run("Oversized Upload Exceeding 55MB is Rejected", func(t *testing.T) {
		boundary := "----CustomBoundary12345"
		header := fmt.Sprintf("--%s\r\nContent-Disposition: form-data; name=\"file\"; filename=\"large.dat\"\r\nContent-Type: application/octet-stream\r\n\r\n", boundary)
		footer := fmt.Sprintf("\r\n--%s--\r\n", boundary)

		// 56 MB stream payload (exceeds 55 MB max body size)
		oversizedPayload := &zeroReader{remaining: 56 * 1024 * 1024}
		bodyReader := io.MultiReader(
			bytes.NewReader([]byte(header)),
			oversizedPayload,
			bytes.NewReader([]byte(footer)),
		)

		req := httptest.NewRequest(http.MethodPost, "/api/upload", bodyReader)
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)

		rec := httptest.NewRecorder()
		h.UploadFile(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 Bad Request, got %d", rec.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode json response: %v", err)
		}

		if code, _ := resp["code"].(string); code != "FILE_TOO_LARGE" {
			t.Errorf("expected error code 'FILE_TOO_LARGE', got %q (resp: %v)", code, resp)
		}
	})

	t.Run("Valid Upload Under Limit Succeeds", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "valid.png")
		if err != nil {
			t.Fatalf("create form file: %v", err)
		}
		_, _ = part.Write([]byte("image content"))
		_ = writer.WriteField("category", "avatar")
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		h.UploadFile(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201 Created, got %d (body: %s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("Unauthorized User Rejected", func(t *testing.T) {
		unauthMw := &mockAuthMiddleware{userID: 0}
		unauthH := uploadHandler.NewUploadHandler(uploadSvc, unauthMw, logger)

		req := httptest.NewRequest(http.MethodPost, "/api/upload", nil)
		rec := httptest.NewRecorder()
		unauthH.UploadFile(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401 Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("UploadBatch Success", func(t *testing.T) {
		uploadSvc := &mockUploadService{
			batchResults: []domain.BatchUploadResultItem{
				{
					OriginalName: "img1.png",
					FileURL:      "/static/uploads/lesson_media/2026/10/uuid1.png",
					SizeBytes:    100,
					MimeType:     "image/png",
				},
				{
					OriginalName: "img2.jpg",
					FileURL:      "/static/uploads/lesson_media/2026/10/uuid2.jpg",
					SizeBytes:    200,
					MimeType:     "image/jpeg",
				},
			},
		}
		authMw := &mockAuthMiddleware{userID: 42}
		h := uploadHandler.NewUploadHandler(uploadSvc, authMw, logger)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part1, err := writer.CreateFormFile("files[]", "img1.png")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		_, _ = part1.Write([]byte("fake png content"))

		part2, err := writer.CreateFormFile("files[]", "img2.jpg")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		_, _ = part2.Write([]byte("fake jpg content"))

		_ = writer.WriteField("category", "lesson_media")
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/upload/batch", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		h.UploadBatch(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201 Created, got %d (body: %s)", rec.Code, rec.Body.String())
		}

		var resp dto.BatchImageUploadResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(resp.Uploaded) != 2 {
			t.Fatalf("expected 2 uploaded items, got %d", len(resp.Uploaded))
		}
		if resp.Uploaded[0].OriginalName != "img1.png" || resp.Uploaded[0].FileURL != "/static/uploads/lesson_media/2026/10/uuid1.png" {
			t.Errorf("unexpected item 0: %+v", resp.Uploaded[0])
		}
	})

	t.Run("UploadBatch Unauthorized", func(t *testing.T) {
		unauthMw := &mockAuthMiddleware{userID: 0}
		unauthH := uploadHandler.NewUploadHandler(uploadSvc, unauthMw, logger)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/upload/batch", nil)
		rec := httptest.NewRecorder()
		unauthH.UploadBatch(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401 Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("UploadBatch Empty Files Rejected", func(t *testing.T) {
		authMw := &mockAuthMiddleware{userID: 42}
		h := uploadHandler.NewUploadHandler(uploadSvc, authMw, logger)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("category", "lesson_media")
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/upload/batch", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		h.UploadBatch(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 Bad Request, got %d", rec.Code)
		}
	})
}

func TestStaticFileServer_PDFHeadersAndRangeRequests(t *testing.T) {
	tempDir := t.TempDir()

	pdfContent := []byte("%PDF-1.4 header bytes and body content for presentation testing")
	pdfPath := filepath.Join(tempDir, "presentation.pdf")
	if err := os.WriteFile(pdfPath, pdfContent, 0o600); err != nil {
		t.Fatalf("failed to write test pdf: %v", err)
	}

	txtContent := []byte("plain text content")
	txtPath := filepath.Join(tempDir, "sample.txt")
	if err := os.WriteFile(txtPath, txtContent, 0o600); err != nil {
		t.Fatalf("failed to write test txt: %v", err)
	}

	server := uploadHandler.StaticFileServer(http.Dir(tempDir))

	t.Run("PDF File has correct Content-Type, Content-Disposition and Accept-Ranges", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/presentation.pdf", nil)
		rec := httptest.NewRecorder()

		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200 OK, got %d", rec.Code)
		}

		if got := rec.Header().Get("Content-Type"); got != "application/pdf" {
			t.Errorf("expected Content-Type application/pdf, got %q", got)
		}

		if got := rec.Header().Get("Content-Disposition"); got != `inline; filename="presentation.pdf"` {
			t.Errorf("expected Content-Disposition inline; filename=\"presentation.pdf\", got %q", got)
		}

		if got := rec.Header().Get("Accept-Ranges"); got != "bytes" {
			t.Errorf("expected Accept-Ranges bytes, got %q", got)
		}

		if !bytes.Equal(rec.Body.Bytes(), pdfContent) {
			t.Errorf("expected body to match pdfContent")
		}
	})

	t.Run("PDF File supports Range requests with 206 Partial Content", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/presentation.pdf", nil)
		req.Header.Set("Range", "bytes=0-7")
		rec := httptest.NewRecorder()

		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusPartialContent {
			t.Fatalf("expected status 206 Partial Content, got %d", rec.Code)
		}

		if rec.Header().Get("Content-Range") == "" {
			t.Errorf("expected Content-Range header present")
		}

		expectedPart := pdfContent[:8]
		if !bytes.Equal(rec.Body.Bytes(), expectedPart) {
			t.Errorf("expected body %q, got %q", string(expectedPart), rec.Body.String())
		}
	})

	t.Run("Non-PDF File does not have presentation Content-Disposition", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sample.txt", nil)
		rec := httptest.NewRecorder()

		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200 OK, got %d", rec.Code)
		}

		if got := rec.Header().Get("Content-Disposition"); got != "" {
			t.Errorf("expected no Content-Disposition for txt file, got %q", got)
		}
	})
}

