package courses

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"edtech/internal/dto"
	errorsAPP "edtech/pkg/errors"
)

func parseQueryParamPagination(r *http.Request) dto.PaginationRequest {
	pagination := dto.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}
	query := r.URL.Query()

	if pageStr := query.Get(string(dto.ParamPage)); pageStr != "" {
		if page, err := strconv.ParseInt(pageStr, 10, 64); err == nil && page > 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := query.Get(string(dto.ParamPageSize)); pageSizeStr != "" {
		if pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64); err == nil && pageSize > 0 {
			pagination.PageSize = pageSize
			if pagination.PageSize > 100 {
				pagination.PageSize = 100
			}
		}
	} else if limitStr := query.Get(string(dto.ParamLimit)); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil && limit > 0 {
			pagination.PageSize = limit
			if pagination.PageSize > 100 {
				pagination.PageSize = 100
			}
		}
	}

	return pagination
}

func parseQueryParamCourseFilter(r *http.Request) (dto.CourseFilterRequest, error) {
	var filter dto.CourseFilterRequest
	query := r.URL.Query()

	if search := query.Get(string(dto.ParamSearch)); search != "" {
		filter.Search = &search
	}

	var err error
	if filter.CategoryID, err = parseQueryParamInt64(dto.ParamCategoryID, &query); err != nil {
		return filter, err
	}

	if filter.CreatedBy, err = parseQueryParamInt64(dto.ParamCreatedBy, &query); err != nil {
		return filter, err
	}

	if diff := query.Get(string(dto.ParamDifficulty)); diff != "" {
		filter.Difficulty = &diff
	}

	if lang := query.Get(string(dto.ParamLanguage)); lang != "" {
		filter.Language = &lang
	}

	filter.SortBy = query.Get(string(dto.ParamSortBy))
	filter.SortOrder = query.Get(string(dto.ParamSortOrder))

	return filter, nil
}

func parseQueryParamInt64(param dto.ParamName, query *url.Values) (*int64, error) {
	if valStr := query.Get(string(param)); valStr != "" {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s must be a valid integer", errorsAPP.ErrValidationFailed, string(param))
		}
		return &val, nil
	}
	return nil, nil
}
