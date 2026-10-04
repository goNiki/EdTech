package upload

import (
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/interfaces/response"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"
)

type UploadHandler struct {
	uploadService  service.UploadServices
	authMiddleware auth.AuthMiddleware
	log            *slog.Logger
}

func NewUploadHandler(
	uploadService service.UploadServices,
	authMiddleware auth.AuthMiddleware,
	log *slog.Logger,
) *UploadHandler {
	return &UploadHandler{
		uploadService:  uploadService,
		authMiddleware: authMiddleware,
		log:            log,
	}
}

// StaticFileServer returns an http.Handler that wraps http.FileServer to serve static files
// with specialized headers for inline viewing and Range streaming (e.g. for PDF presentations).
func StaticFileServer(root http.FileSystem) http.Handler {
	fileServer := http.FileServer(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(strings.ToLower(r.URL.Path), ".pdf") {
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", `inline; filename="presentation.pdf"`)
			w.Header().Set("Accept-Ranges", "bytes")
		}
		fileServer.ServeHTTP(w, r)
	})
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.upload.UploadFile"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	// 55 MB max upload request body size to accommodate presentations up to 50 MB
	const maxUploadBodySize = 55 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBodySize)

	// 16 MB in-memory parsing buffer; excess is stored in temp files or rejected by MaxBytesReader
	if err := r.ParseMultipartForm(16 << 20); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			response.HandleError(w, r, h.log, errorsAPP.ErrFileTooLarge, op)
			return
		}
		response.HandleError(w, r, h.log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrValidationFailed, op)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	category := r.FormValue("category")

	result, err := h.uploadService.UploadFile(r.Context(), file, header.Filename, header.Size, category)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	resp := dto.FileUploadResponse{
		FileURL:   result.FileURL,
		FileName:  result.FileName,
		SizeBytes: result.SizeBytes,
		MimeType:  result.MimeType,
	}

	response.Created(w, r, resp)
}

func (h *UploadHandler) UploadBatch(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.upload.UploadBatch"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	// 55 MB max batch upload request body size to protect against disk exhaustion DoS
	const maxBatchUploadBodySize = 55 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxBatchUploadBodySize)

	// 32 MB in-memory parsing buffer; excess is stored in temp files or rejected by MaxBytesReader
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			response.HandleError(w, r, h.log, errorsAPP.ErrFileTooLarge, op)
			return
		}
		response.HandleError(w, r, h.log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrEmptyFile, op)
		return
	}

	fileHeaders := r.MultipartForm.File["files[]"]
	if len(fileHeaders) == 0 {
		fileHeaders = r.MultipartForm.File["files"]
	}
	if len(fileHeaders) == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrEmptyFile, op)
		return
	}

	category := r.FormValue("category")
	if category == "" {
		category = "lesson_media"
	}

	batchItems := make([]domain.BatchFileItem, 0, len(fileHeaders))
	openFiles := make([]multipart.File, 0, len(fileHeaders))
	defer func() {
		for _, f := range openFiles {
			_ = f.Close()
		}
	}()

	for _, fh := range fileHeaders {
		f, err := fh.Open()
		if err != nil {
			response.HandleError(w, r, h.log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
			return
		}
		openFiles = append(openFiles, f)
		batchItems = append(batchItems, domain.BatchFileItem{
			Reader:   f,
			Filename: fh.Filename,
			Size:     fh.Size,
		})
	}

	results, err := h.uploadService.UploadImagesBatch(r.Context(), batchItems, category)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	uploadedItems := make([]dto.BatchImageItemResponse, len(results))
	for i, res := range results {
		uploadedItems[i] = dto.BatchImageItemResponse{
			OriginalName: res.OriginalName,
			FileURL:      res.FileURL,
			SizeBytes:    res.SizeBytes,
		}
	}

	response.Created(w, r, dto.BatchImageUploadResponse{
		Uploaded: uploadedItems,
	})
}

