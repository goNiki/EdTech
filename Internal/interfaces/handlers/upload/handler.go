package upload

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

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

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.upload.UploadFile"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	// 30 MB max upload request body size to protect against disk exhaustion DoS
	const maxUploadBodySize = 30 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBodySize)

	// 10 MB in-memory parsing buffer; excess is rejected by MaxBytesReader
	if err := r.ParseMultipartForm(10 << 20); err != nil {
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
