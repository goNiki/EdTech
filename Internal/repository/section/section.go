package section

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository/models"
	"edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repositorySection) ListSectionsByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64) ([]domain.Section, error) {
	const op = "repository.section.ListSectionsByCourseID"

	query := `
		SELECT id, course_id, title, description, position, created_at, updated_at, deleted_at 
		FROM sections 
		WHERE course_id = $1 AND deleted_at IS NULL 
		ORDER BY position ASC
	`

	rows, err := q.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var sections []domain.Section
	for rows.Next() {
		var m models.Section
		if err := rows.Scan(
			&m.ID,
			&m.CourseID,
			&m.Title,
			&m.Description,
			&m.Position,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan section: %w", op, err)
		}
		sections = append(sections, *converter.SectionToDomain(&m))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w", op, err)
	}

	if sections == nil {
		sections = []domain.Section{}
	}

	return sections, nil
}

func (r *repositorySection) GetSectionByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Section, error) {
	const op = "repository.section.GetSectionByID"

	query := `
		SELECT id, course_id, title, description, position, created_at, updated_at, deleted_at 
		FROM sections 
		WHERE id = $1 AND deleted_at IS NULL
	`

	var m models.Section
	err := q.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.CourseID,
		&m.Title,
		&m.Description,
		&m.Position,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrValidationFailed)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return converter.SectionToDomain(&m), nil
}

func (r *repositorySection) CreateSection(ctx context.Context, q db.QueryExecutor, section *domain.Section) (*domain.Section, error) {
	const op = "repository.section.CreateSection"

	entity := converter.SectionToEntity(section)

	query := `
		INSERT INTO sections (course_id, title, description, position)
		VALUES ($1, $2, $3, $4)
		RETURNING id, course_id, title, description, position, created_at, updated_at, deleted_at
	`

	err := q.QueryRow(ctx, query,
		entity.CourseID,
		entity.Title,
		entity.Description,
		entity.Position,
	).Scan(
		&entity.ID,
		&entity.CourseID,
		&entity.Title,
		&entity.Description,
		&entity.Position,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return converter.SectionToDomain(entity), nil
}

func (r *repositorySection) UpdateSection(ctx context.Context, q db.QueryExecutor, section *domain.Section) error {
	const op = "repository.section.UpdateSection"

	query := `
		UPDATE sections 
		SET title = $1, description = $2, position = $3, updated_at = NOW() 
		WHERE id = $4 AND deleted_at IS NULL
	`

	tag, err := q.Exec(ctx, query, section.Title, section.Description, section.Position, section.ID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}

func (r *repositorySection) DeleteSection(ctx context.Context, q db.QueryExecutor, sectionID int64) error {
	const op = "repository.section.DeleteSection"

	query := `
		UPDATE sections 
		SET deleted_at = NOW(), updated_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := q.Exec(ctx, query, sectionID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}

func (r *repositorySection) GetMaxPositionByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64) (int, error) {
	const op = "repository.section.GetMaxPositionByCourseID"

	query := `SELECT COALESCE(MAX(position), 0) FROM sections WHERE course_id = $1 AND deleted_at IS NULL`

	var maxPos int
	err := q.QueryRow(ctx, query, courseID).Scan(&maxPos)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return maxPos, nil
}
