package courses

import (
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *CourseHandler) GetCourseByID(w http.ResponseWriter, r *http.Request) {
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

	permissionsDomain, err := h.accessService.BuildCoursePermissions(r.Context(), course, userID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	userRole := ""
	role, err := h.enrolmentService.GetRoleUserInCource(r.Context(), userID, course.Id)
	if err != nil {
		if err.Error() != "user not found" && err.Error() != "not found" && err.Error() != "no rows in result set" {
			response.HandleError(w, r, log, err, op)
			return
		}
	} else {
		userRole = role
	}

	resp := dto.CourseDetailResponse{
		Course:      converter.CourseToDTO(course),
		Permissions: converter.CoursePermissionsToDTO(permissionsDomain),
		UserRole:    userRole,
	}

	response.OK(w, r, resp)
}

func (h *CourseHandler) GetCourseBySlug(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.GetCourseBySlug"

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

	permissionsDomain, err := h.accessService.BuildCoursePermissions(r.Context(), course, userID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	userRole := ""
	role, err := h.enrolmentService.GetRoleUserInCource(r.Context(), userID, course.Id)
	if err != nil {
		if err.Error() != "user not found" && err.Error() != "not found" && err.Error() != "no rows in result set" {
			response.HandleError(w, r, log, err, op)
			return
		}
	} else {
		userRole = role
	}

	resp := dto.CourseDetailResponse{
		Course:      converter.CourseToDTO(course),
		Permissions: converter.CoursePermissionsToDTO(permissionsDomain),
		UserRole:    userRole,
	}

	response.OK(w, r, resp)
}

func (h *CourseHandler) GetCourseStructure(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.GetCourseStructure"

	log := logger.GetLogger(r.Context(), op)

	courseIDStr := chi.URLParam(r, "courseid")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	result, err := h.courseService.GetCourseStructure(r.Context(), courseID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, result)
}
