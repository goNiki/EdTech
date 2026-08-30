package courses

import (
	"fmt"
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// ListMyCourses handles GET /api/v1/courses/my with query filters, role filtering, search, sorting and pagination
func (h *CourseHandler) ListMyCourses(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.ListMyCourses"

	log := logger.GetLogger(r.Context(), op)

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	var req dto.ListMyCoursesRequest
	var err error

	req.UserID = userID
	req.Role = r.URL.Query().Get(string(dto.ParamRole))
	req.Pagination = parseQueryParamPagination(r)

	req.Filter, err = parseQueryParamCourseFilter(r)
	if err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	inputDomain := converter.ListMyCoursesRequestToDomain(req)

	paginated, err := h.courseService.ListMyCourses(r.Context(), inputDomain)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	dtoCourses := make([]dto.Course, 0, len(paginated.Courses))
	for _, c := range paginated.Courses {
		dtoCourses = append(dtoCourses, converter.CourseToDTO(&c))
	}

	respData := dto.PaginatedCourses{
		Courses:  dtoCourses,
		Page:     paginated.Page,
		PageSize: paginated.PageSize,
		Total:    paginated.Total,
	}

	response.OK(w, r, respData)
}
