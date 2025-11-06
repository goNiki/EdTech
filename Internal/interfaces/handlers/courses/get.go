package courses

import (
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/response"
	errorsAPP "edtech/pkg/errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

	courseID := chi.URLParam(r, "courseid")
	courseIDInt, err := strconv.Atoi(courseID)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}
	userID := h.authMiddleware.GetUserID(r.Context())

	course, err := h.courseService.GetCourseByID(r.Context(), int64(courseIDInt))
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	canView, err := h.accessService.CanViewCourse(r.Context(), course, userID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}
	if !canView {
		response.HandleError(w, r, log, errorsAPP.ErrForbidden, op)
		return
	}

	permission := h.BuildCoursePermissions(r.Context(), course, userID)

	userRole, err := h.enrolmentService.GetRoleUserInCource(r.Context(), userID, course.Id)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	date := dto.Course{
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

	resp := dto.CourseDetailResponse{
		Course:      date,
		Permissions: permission,
		UserRole:    userRole,
	}

	response.OK(w, r, resp)

}

func (h *handler) GetCourseBySlug(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.get"

	log := logger.GetLogger(r.Context(), op)

	courseSlug := chi.URLParam(r, "slug")

	course, err := h.courseService.GetCourseBySlug(r.Context(), courseSlug)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}
	userID := h.authMiddleware.GetUserID(r.Context())

	canView, err := h.accessService.CanViewCourse(r.Context(), course, userID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	if !canView {
		response.HandleError(w, r, log, errorsAPP.ErrForbidden, op)
		return
	}

	permission := h.BuildCoursePermissions(r.Context(), course, userID)

	userRole, err := h.enrolmentService.GetRoleUserInCource(r.Context(), userID, course.Id)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	date := dto.Course{
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

	resp := dto.CourseDetailResponse{
		Course:      date,
		Permissions: permission,
		UserRole:    userRole,
	}

	response.OK(w, r, resp)

}
