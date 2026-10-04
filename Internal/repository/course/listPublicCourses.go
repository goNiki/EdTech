package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository/models"
	"edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"
)

func buildCourseFilterQuery(filter domain.CourseFilter) (string, []any) {
	where := " WHERE status = 'published' AND visibility = 'public' AND deleted_at IS NULL"
	var args []any
	argIdx := 1

	if filter.CategoryID != nil {
		where += fmt.Sprintf(" AND category_id = $%d", argIdx)
		args = append(args, *filter.CategoryID)
		argIdx++
	}

	if filter.CreatedBy != nil {
		where += fmt.Sprintf(" AND created_by = $%d", argIdx)
		args = append(args, *filter.CreatedBy)
		argIdx++
	}

	if filter.Difficulty != nil && *filter.Difficulty != "" {
		where += fmt.Sprintf(" AND difficulty = $%d", argIdx)
		args = append(args, *filter.Difficulty)
		argIdx++
	}

	if filter.Language != nil && *filter.Language != "" {
		where += fmt.Sprintf(" AND language = $%d", argIdx)
		args = append(args, *filter.Language)
		argIdx++
	}

	if filter.Search != nil && *filter.Search != "" {
		where += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+*filter.Search+"%")
	}

	return where, args
}

func (r *repository) ListPublicCourses(ctx context.Context, pagination domain.Pagination, filter domain.CourseFilter) ([]domain.Course, error) {
	const op = "repository.course.ListPublicCourses"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	whereClause, args := buildCourseFilterQuery(filter)

	// Whitelist безопасной сортировки (защита от SQL-инъекций)
	sortCol := "created_at"
	switch filter.SortBy {
	case "title", "enrolled_count", "created_at", "total_lessons":
		sortCol = filter.SortBy
	}

	sortDir := "DESC"
	if filter.SortOrder == "asc" || filter.SortOrder == "ASC" {
		sortDir = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT 
			id, title, slug, short_description, description, cover_url, intro_video_url, 
			created_by, visibility, status, difficulty, language, estimated_duration, 
			category_id, total_lessons, total_sections, enrolled_count, 
			created_at, updated_at, published_at, archived_at, deleted_at 
		FROM courses 
		%s
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d
	`, whereClause, sortCol, sortDir, len(args)+1, len(args)+2)

	args = append(args, pagination.PageSize, pagination.Offset())

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var m models.Course
		if err := rows.Scan(
			&m.ID, &m.Title, &m.Slug, &m.ShortDescription, &m.Description, &m.CoverURL, &m.IntroVideoURL,
			&m.CreatedBy, &m.Visibility, &m.Status, &m.Difficulty, &m.Language, &m.EstimatedDuration,
			&m.CategoryID, &m.TotalLessons, &m.TotalSections, &m.EnrolledCount,
			&m.CreatedAt, &m.UpdatedAt, &m.PublishedAt, &m.ArchivedAt, &m.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan course: %w", op, err)
		}
		courses = append(courses, *converter.CourseToDomain(&m))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w", op, err)
	}

	if courses == nil {
		courses = []domain.Course{}
	}

	return courses, nil
}
