package course

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) UpdateCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error {
	const op = "repository.course.updatecourse"

	query := `
		UPDATE courses
		SET 
			title = $1, 
			slug = $2, 
			short_description = $3, 
			description = $4, 
			cover_url = $5, 
			intro_video_url = $6, 
			visibility = $7, 
			status = $8, 
			difficulty = $9, 
			language = $10, 
			estimated_duration = $11, 
			category_id = $12,
			updated_at = NOW()
		WHERE id = $13 AND deleted_at IS NULL AND updated_at = $14
	`

	entityCourse := repoconverter.CourseToEntity(course)

	cmtTag, err := q.Exec(
		ctx,
		query,
		entityCourse.Title,
		entityCourse.Slug,
		entityCourse.ShortDescription,
		entityCourse.Description,
		entityCourse.CoverURL,
		entityCourse.IntroVideoURL,
		entityCourse.Visibility,
		entityCourse.Status,
		entityCourse.Difficulty,
		entityCourse.Language,
		entityCourse.EstimatedDuration,
		entityCourse.CategoryID,
		entityCourse.ID,
		entityCourse.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, errorsAPP.ErrSlugAlreadyExists)
		}
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if cmtTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}
