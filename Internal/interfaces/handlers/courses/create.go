package courses

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
	"net/http"

	"github.com/go-chi/render"
)

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.create"

	log := logger.GetLogger(r.Context(), op)

	var req dto.CreateCourseRequest
	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}
	userRole := h.authMiddleware.GetUserRole(r.Context())
	if userRole != string(domain.RoleTeacher) {
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

	createdCourse, err := h.courseService.CreateCourse(r.Context(), &course)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	resp := dto.CreateCourseResponse{
		Data: dto.CreateCourseResponseData{
			ID:        createdCourse.Id,
			Title:     createdCourse.Title,
			Slug:      createdCourse.Slug,
			Status:    createdCourse.Status,
			CreatedAt: createdCourse.CreatedAt,
		},
	}
	response.Created(w, r, resp.Data)
}
