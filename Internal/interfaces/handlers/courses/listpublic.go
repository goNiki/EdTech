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

// ListPublicCourses handles GET /api/v1/courses with query filters, search, sorting and pagination
func (h *CourseHandler) ListPublicCourses(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.ListPublicCourses"

	log := logger.GetLogger(r.Context(), op)

	var req dto.ListPublicCoursesRequest
	var err error

	req.Pagination = parseQueryParamPagination(r)

	req.Filter, err = parseQueryParamCourseFilter(r)
	if err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	paginationDomain, filterDomain := converter.ListPublicCoursesRequestToDomain(req)

	paginated, err := h.courseService.ListPublicCourses(r.Context(), paginationDomain, filterDomain)
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
