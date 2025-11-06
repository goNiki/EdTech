package courses

import (
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/response"
	errorsAPP "edtech/pkg/errors"
	"net/http"
	"strconv"
)

// page (default: 1)
// page_size (default: 10, max: 100)
func (h *handler) ListCourses(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.list"

	log := logger.GetLogger(r.Context(), op)

	page := r.URL.Query().Get("page")

	pagesize := r.URL.Query().Get("pagesize")

	if page == "" {
		page = "1"
	}
	if pagesize == "" {
		pagesize = "10"
	}
	if pagesize > "100" || pagesize < "1" {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLQuery, op)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}
	pagesizeInt, err := strconv.Atoi(pagesize)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	courses, err := h.courseService.ListCourses(r.Context(), int64(pageInt), int64(pagesizeInt))
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	respData := dto.PaginatedCourses{
		Courses:  courses.Courses,
		Page:     courses.Page,
		PageSize: courses.PageSize,
		Total:    courses.Total,
	}

	response.OK(w, r, respData)
}
