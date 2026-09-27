package section

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/interfaces/handlers/converter"
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
}

func NewSectionHandler(sectionService service.SectionServices, log *slog.Logger, validator *validator.Validate) *SectionHandler {
	return &SectionHandler{
		sectionService: sectionService,
		log:            log,
		validator:      validator,
	}
}

func (h *SectionHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.section.CreateSection"

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

	createdSection, err := h.sectionService.CreateSection(r.Context(), &sectionDomain)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.Created(w, r, converter.SectionToDTO(createdSection))
}

func (h *SectionHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.section.UpdateSection"

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

	if err := h.sectionService.UpdateSection(r.Context(), sectionDomain); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"status": "success"})
}

func (h *SectionHandler) DeleteSection(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.section.DeleteSection"

	idStr := chi.URLParam(r, "id")
	sectionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	if err := h.sectionService.DeleteSection(r.Context(), sectionID); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"status": "deleted"})
}
