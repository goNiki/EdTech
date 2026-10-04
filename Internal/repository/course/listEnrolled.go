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

func buildEnrolledFilterQuery(input *domain.InputListMyCourse) (string, []any) {
	where := " WHERE uc.user_id = $1 AND c.deleted_at IS NULL"
	args := []any{input.UserID}
	argIdx := 2

	if input.Role == string(domain.StudentRole) {
		where += fmt.Sprintf(" AND uc.role = $%d AND c.status = 'published'", argIdx)
		args = append(args, "student")
		argIdx++
	} else if input.Role == string(domain.TeacherRole) || input.Role == string(domain.CreatorRole) {
		where += fmt.Sprintf(" AND uc.role IN ($%d, $%d)", argIdx, argIdx+1)
		args = append(args, "teacher", "creator")
		argIdx += 2
	} else if input.Role != "" {
		where += fmt.Sprintf(" AND uc.role = $%d", argIdx)
		args = append(args, input.Role)
		argIdx++
	}

	if input.Filter.CategoryID != nil {
		where += fmt.Sprintf(" AND c.category_id = $%d", argIdx)
		args = append(args, *input.Filter.CategoryID)
		argIdx++
	}

	if input.Filter.Difficulty != nil && *input.Filter.Difficulty != "" {
		where += fmt.Sprintf(" AND c.difficulty = $%d", argIdx)
		args = append(args, *input.Filter.Difficulty)
		argIdx++
	}

	if input.Filter.Language != nil && *input.Filter.Language != "" {
		where += fmt.Sprintf(" AND c.language = $%d", argIdx)
		args = append(args, *input.Filter.Language)
		argIdx++
	}

	if input.Filter.Search != nil && *input.Filter.Search != "" {
		where += fmt.Sprintf(" AND (c.title ILIKE $%d OR c.description ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+*input.Filter.Search+"%")
	}

	return where, args
}

func (r *repository) ListEnrolledCourses(ctx context.Context, input *domain.InputListMyCourse) ([]domain.Course, error) {
	const op = "repository.course.ListEnrolledCourses"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	whereClause, args := buildEnrolledFilterQuery(input)

	// Whitelist безопасной сортировки
	sortCol := "uc.enrolled_at"
	switch input.Filter.SortBy {
	case "created_at":
		sortCol = "c.created_at"
	case "title":
		sortCol = "c.title"
	case "enrolled_count":
		sortCol = "c.enrolled_count"
	case "enrolled_at":
		sortCol = "uc.enrolled_at"
	}

	sortDir := "DESC"
	if input.Filter.SortOrder == "asc" || input.Filter.SortOrder == "ASC" {
		sortDir = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT 
			c.id, c.title, c.slug, c.short_description, c.description, c.cover_url, c.intro_video_url, 
			c.created_by, c.visibility, c.status, c.difficulty, c.language, c.estimated_duration, 
			c.category_id, c.total_lessons, c.total_sections, c.enrolled_count, 
			c.created_at, c.updated_at, c.published_at, c.archived_at, c.deleted_at 
		FROM courses c
		INNER JOIN users_courses uc ON c.id = uc.course_id 
		%s
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d
	`, whereClause, sortCol, sortDir, len(args)+1, len(args)+2)

	args = append(args, input.Pagination.PageSize, input.Pagination.Offset())

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
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	if courses == nil {
		courses = []domain.Course{}
	}

	return courses, nil
}

func (r *repository) CountEnrolledCourses(ctx context.Context, input *domain.InputListMyCourse) (int64, error) {
	const op = "repository.course.CountEnrolledCourses"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	whereClause, args := buildEnrolledFilterQuery(input)
	query := `
		SELECT COUNT(*) 
		FROM courses c
		INNER JOIN users_courses uc ON c.id = uc.course_id 
	` + whereClause

	var total int64
	err := q.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return total, nil
}
