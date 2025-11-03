package courses

import (
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Получение курса по ID
// Метод: GET /api/v1/courses/{id}
// Доступ: Public (только published), Enrolled users (все)
// Handler: GetCourseByID(w http.ResponseWriter, r *http.Request)
// Response: 200 OK или 404 Not Found
// 2.4 Получение курса по slug
// Метод: GET /api/v1/courses/slug/{slug}
// Доступ: Public
// Handler: GetCourseBySlug(w http.ResponseWriter, r *http.Request)

func (h *handler) GetCourseByID(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.get"

	log := logger.GetLogger(r.Context(), op)

	courseID := chi.URLParam(r, "id")
	courseIDInt, err := strconv.Atoi(courseID)
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	course, err := h.courseService.GetCourseByID(r.Context(), int64(courseIDInt))
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundCourse) {
			log.Error("course not found", sl.Error(err))
			http.Error(w, "course not found", http.StatusNotFound)
			return
		}
		log.Error("internal server error", sl.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.Course{
		ID:          course.Id,
		Title:       course.Title,
		Slug:        course.Slug,
		Description: course.Description,
		CoverURL:    course.CoverURL,
		CreatedBy:   course.CreatedBy,
		Visibility:  course.Visibility,
		Status:      course.Status,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)

}

func (h *handler) GetCourseBySlug(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.get"

	log := logger.GetLogger(r.Context(), op)

	courseSlug := chi.URLParam(r, "slug")

	course, err := h.courseService.GetCourseBySlug(r.Context(), courseSlug)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundCourse) {
			log.Error("course not found", sl.Error(err))
			http.Error(w, "course not found", http.StatusNotFound)
			return
		}
		log.Error("internal server error", sl.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	userID := h.authMiddleware.GetUserID(r.Context())

	if course.Visibility == "private" && course.Status == "draft" && course.CreatedBy != int64(h.authMiddleware.GetUserID(r.Context())) {
		log.Error("course is not published", sl.Error(err))
		http.Error(w, "course is not published", http.StatusForbidden)
		return
	}

	resp := dto.Course{
		ID:          course.Id,
		Title:       course.Title,
		Slug:        course.Slug,
		Description: course.Description,
		CoverURL:    course.CoverURL,
		CreatedBy:   course.CreatedBy,
		Visibility:  course.Visibility,
		Status:      course.Status,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)

}
