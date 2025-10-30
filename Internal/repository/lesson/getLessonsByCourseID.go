package lesson

import (
	"context"
	"edtech/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetLessonsByCourseID(ctx context.Context, courseID int) ([]domain.Lesson, error) {

	const op = "repositiry.couse.lessonrepo.ListLessonsByCourseID"

	query := `SELECT id, course_id, title, description, cover_url, position, created_at, updated_at FROM lessons WHERE course_id = $1 ORDER BY position `

	rows, err := r.Pool.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	var lessons []domain.Lesson

	for rows.Next() {
		var l domain.Lesson

		if err := rows.Scan(&l.ID, &l.CourseID, &l.Title, &l.Description, &l.CoverURL, &l.Position, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		lessons = append(lessons, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(lessons) == 0 {
		return nil, fmt.Errorf("%s: %w", op, pgx.ErrNoRows)
	}

	return lessons, nil

}
