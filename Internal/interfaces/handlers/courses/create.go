package courses

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/response"
	errorsAPP "edtech/pkg/errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

func (h *handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.create"

	log := logger.GetLogger(r.Context(), op)

	var req dto.CreateCourseRequest
	userID := h.authMiddleware.GetUserID(r.Context())
	userRole := h.authMiddleware.GetUserRole(r.Context())
	if userRole != "teacher" {
		response.HandleError(w, r, log, errorsAPP.ErrForbidden, op)
		return
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrCourseValidation, op)
		return
	}

	course := domain.Course{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		CreatedBy:   userID,
		Visibility:  req.Visibility,
	}

	courseID, err := h.courseService.CreateCourse(r.Context(), &course)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	resp := dto.CreateCourseResponse{
		Data: dto.CreateCourseResponseData{
			ID:        courseID,
			Title:     req.Title,
			Slug:      req.Slug,
			Status:    req.Status,
			CreatedAt: time.Now(),
		},
		Message: "Курс успешно создан",
	}
	response.Created(w, r, resp.Data, resp.Message)
}
