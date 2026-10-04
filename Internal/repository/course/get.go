package course

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	const op = "repository.course.GetCourseByID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
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
		FROM courses 
		WHERE id = $1 AND deleted_at IS NULL`

	var entity repomodels.Course

	err := q.QueryRow(
		ctx,
		query,
		id,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotFound)
		}
		return nil, fmt.Errorf("%s: %w :%w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.CourseToDomain(&entity), nil
}

func (r *repository) GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error) {
	const op = "repository.course.getcoursebyslug"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
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
		FROM courses 
		WHERE slug = $1 AND deleted_at IS NULL`

	var entity repomodels.Course

	err := q.QueryRow(
		ctx,
		query,
		slug,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotFound)
		}
		return nil, fmt.Errorf("%s: %w :%w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.CourseToDomain(&entity), nil
}
