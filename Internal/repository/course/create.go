package course

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	const op = "repository.course.CreateCourse"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	entity := converter.CourseToEntity(course)

	query := `
		INSERT INTO courses (
			title, 
			slug, 
			short_description, 
			description, 
			cover_url, 
			intro_video_url, 
			created_by, 
			visibility, 
			status, 
			difficulty, 
			language, 
			estimated_duration, 
			category_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) 
		RETURNING 
			id, 
			title, 
			slug, 
			short_description, 
			description, 
			cover_url, 
			intro_video_url, 
			created_by, 
			visibility, 
			status, 
			difficulty, 
			language, 
			estimated_duration, 
			category_id, 
			total_lessons, 
			total_sections, 
			enrolled_count, 
			created_at, 
			updated_at, 
			published_at, 
			archived_at, 
			deleted_at
	`

	err := q.QueryRow(ctx, query,
		entity.Title,
		entity.Slug,
		entity.ShortDescription,
		entity.Description,
		entity.CoverURL,
		entity.IntroVideoURL,
		entity.CreatedBy,
		entity.Visibility,
		entity.Status,
		entity.Difficulty,
		entity.Language,
		entity.EstimatedDuration,
		entity.CategoryID,
	).Scan(
		&entity.ID,
		&entity.Title,
		&entity.Slug,
		&entity.ShortDescription,
		&entity.Description,
		&entity.CoverURL,
		&entity.IntroVideoURL,
		&entity.CreatedBy,
		&entity.Visibility,
		&entity.Status,
		&entity.Difficulty,
		&entity.Language,
		&entity.EstimatedDuration,
		&entity.CategoryID,
		&entity.TotalLessons,
		&entity.TotalSections,
		&entity.EnrolledCount,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.PublishedAt,
		&entity.ArchivedAt,
		&entity.DeletedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrSlugAlreadyExists)
		}
		return nil, fmt.Errorf("%s: insert course: %w", op, err)
	}

	return converter.CourseToDomain(entity), nil
}
