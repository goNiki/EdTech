package courses

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
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
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed decode json", sl.Error(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		log.Error("failed validate fields", sl.Error(err))
		http.Error(w, "failed validate fields", http.StatusBadRequest)
		return
	}

	course := domain.Course{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		CreatedBy:   int(userID),
		Visibility:  req.Visibility,
	}

	courseID, err := h.courseService.CreateCourse(r.Context(), &course)
	if err != nil {
		log.Error("failed create course", sl.Error(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
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

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}
