package section

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/interfaces/handlers/converter"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/interfaces/response"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type SectionHandler struct {
	sectionService service.SectionServices
	log            *slog.Logger
	validator      *validator.Validate
	authMiddleware auth.AuthMiddleware
}

func NewSectionHandler(sectionService service.SectionServices, log *slog.Logger, validator *validator.Validate, authMiddleware auth.AuthMiddleware) *SectionHandler {
	return &SectionHandler{
		sectionService: sectionService,
		log:            log,
		validator:      validator,
		authMiddleware: authMiddleware,
	}
}

func (h *SectionHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.section.CreateSection"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	var req dto.CreateSectionRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrValidationFailed, op)
		return
	}

	sectionDomain := converter.CreateSectionRequestToDomain(req)

	createdSection, err := h.sectionService.CreateSection(r.Context(), userID, &sectionDomain)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.Created(w, r, converter.SectionToDTO(createdSection))
}

func (h *SectionHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.section.UpdateSection"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	idStr := chi.URLParam(r, "id")
	sectionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.UpdateSectionRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrValidationFailed, op)
		return
	}

	sectionDomain := &domain.Section{
		ID: sectionID,
	}

	if req.Title != nil {
		sectionDomain.Title = *req.Title
	}
	if req.Description != nil {
		sectionDomain.Description = *req.Description
	}

	if err := h.sectionService.UpdateSection(r.Context(), userID, sectionDomain); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"status": "success"})
}

func (h *SectionHandler) DeleteSection(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.section.DeleteSection"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	idStr := chi.URLParam(r, "id")
	sectionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	if err := h.sectionService.DeleteSection(r.Context(), userID, sectionID); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"status": "deleted"})
}
