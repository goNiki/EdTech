package upload

import (
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

	// 32 MB max memory in form parsing
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrFileTooLarge, op)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrValidationFailed, op)
		return
	}
	defer file.Close()

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
