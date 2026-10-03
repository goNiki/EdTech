package lesson

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

type LessonHandler struct {
	lessonService  service.LessonServices
	log            *slog.Logger
	validator      *validator.Validate
	authMiddleware auth.AuthMiddleware
}

func NewLessonHandler(lessonService service.LessonServices, log *slog.Logger, validator *validator.Validate, authMiddleware auth.AuthMiddleware) *LessonHandler {
	return &LessonHandler{
		lessonService:  lessonService,
		log:            log,
		validator:      validator,
		authMiddleware: authMiddleware,
	}
}

func (h *LessonHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.lesson.CreateLesson"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	var req dto.CreateLessonRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrValidationFailed, op)
		return
	}

	lessonDomain := converter.CreateLessonRequestToDomain(req)

	createdID, err := h.lessonService.CreateLesson(r.Context(), userID, &lessonDomain)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.Created(w, r, map[string]int64{"id": createdID})
}

func (h *LessonHandler) UpdateLesson(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.lesson.UpdateLesson"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	idStr := chi.URLParam(r, "id")
	lessonID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.UpdateLessonRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrValidationFailed, op)
		return
	}

	lessonDomain := &domain.Lesson{
		ID: lessonID,
	}

	if req.SectionID != nil {
		lessonDomain.SectionID = req.SectionID
	}
	if req.Title != nil {
		lessonDomain.Title = *req.Title
	}
	if req.Description != nil {
		lessonDomain.Description = *req.Description
	}
	if req.CoverURL != nil {
		lessonDomain.CoverURL = *req.CoverURL
	}
	if req.Content != nil {
		lessonDomain.Content = *req.Content
	}
	if req.Type != nil {
		lessonDomain.Type = *req.Type
	}
	if req.Duration != nil {
		lessonDomain.Duration = req.Duration
	}
	if req.IsFree != nil {
		lessonDomain.IsFree = *req.IsFree
	}

	if err := h.lessonService.UpdateLesson(r.Context(), userID, lessonDomain); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"status": "success"})
}

func (h *LessonHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.lesson.DeleteLesson"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	idStr := chi.URLParam(r, "id")
	lessonID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	if err := h.lessonService.DeleteLesson(r.Context(), userID, lessonID); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"status": "deleted"})
}
