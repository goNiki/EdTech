package courses

import (
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

// page (default: 1)
// page_size (default: 10, max: 100)
func (h *handler) ListCourses(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.list"

	log := logger.GetLogger(r.Context(), op)

	page, pagesize := r.URL.Query().Get("page"), r.URL.Query().Get("pagesize")
	if page == "" {
		page = "1"
	}
	if pagesize == "" {
		pagesize = "10"
	}
	if pagesize > "100" || pagesize < "1" {
		http.Error(w, "pagesize is too large", http.StatusBadRequest)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "invalid page", http.StatusBadRequest)
		return
	}
	pagesizeInt, err := strconv.Atoi(pagesize)
	if err != nil {
		http.Error(w, "invalid pagesize", http.StatusBadRequest)
		return
	}

	courses, err := h.courseService.ListCourses(r.Context(), pageInt, pagesizeInt)
	if err != nil {
		log.Error("failed list courses", sl.Error(err))
		http.Error(w, "failed list courses", http.StatusInternalServerError)
		return
	}

	respData := dto.PaginatedCourses{
		Courses:  courses.Courses,
		Page:     pageInt,
		PageSize: pagesizeInt,
		Total:    courses.Total,
	}
	resp := dto.ListCoursesResponse{
		Data: respData,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
