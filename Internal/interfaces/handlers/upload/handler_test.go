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
	"testing"

	"edtech/internal/domain"
	uploadHandler "edtech/internal/interfaces/handlers/upload"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"
)

type mockUploadService struct {
	service.UploadServices
	uploadedFile *domain.FileUploadResult
	err          error
}

func (m *mockUploadService) UploadFile(ctx context.Context, file io.Reader, filename string, size int64, category string) (*domain.FileUploadResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.uploadedFile, nil
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

	t.Run("Oversized Upload Exceeding 30MB is Rejected", func(t *testing.T) {
		boundary := "----CustomBoundary12345"
		header := fmt.Sprintf("--%s\r\nContent-Disposition: form-data; name=\"file\"; filename=\"large.dat\"\r\nContent-Type: application/octet-stream\r\n\r\n", boundary)
		footer := fmt.Sprintf("\r\n--%s--\r\n", boundary)

		// 31 MB stream payload
		oversizedPayload := &zeroReader{remaining: 31 * 1024 * 1024}
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
}
