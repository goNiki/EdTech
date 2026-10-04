package section

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository/models"
	"edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repositorySection) ListSectionsByCourseID(ctx context.Context, courseID int64) ([]domain.Section, error) {
	const op = "repository.section.ListSectionsByCourseID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT id, course_id, title, description, position, status, created_at, updated_at, deleted_at 
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
			&m.Status,
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

func (r *repositorySection) GetSectionByID(ctx context.Context, id int64) (*domain.Section, error) {
	const op = "repository.section.GetSectionByID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT id, course_id, title, description, position, status, created_at, updated_at, deleted_at 
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
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrSectionNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return converter.SectionToDomain(&m), nil
}

func (r *repositorySection) CreateSection(ctx context.Context, section *domain.Section) (*domain.Section, error) {
	const op = "repository.section.CreateSection"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	entity := converter.SectionToEntity(section)
	if entity.Status == "" {
		entity.Status = domain.StatusDraft
	}

	query := `
		INSERT INTO sections (course_id, title, description, position, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, course_id, title, description, position, status, created_at, updated_at, deleted_at
	`

	err := q.QueryRow(ctx, query,
		entity.CourseID,
		entity.Title,
		entity.Description,
		entity.Position,
		entity.Status,
	).Scan(
		&entity.ID,
		&entity.CourseID,
		&entity.Title,
		&entity.Description,
		&entity.Position,
		&entity.Status,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return converter.SectionToDomain(entity), nil
}

func (r *repositorySection) UpdateSection(ctx context.Context, section *domain.Section) error {
	const op = "repository.section.UpdateSection"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		UPDATE sections 
		SET title = $1, description = $2, position = $3, status = $4, updated_at = NOW() 
		WHERE id = $5 AND deleted_at IS NULL
	`

	status := section.Status
	if status == "" {
		status = domain.StatusDraft
	}

	tag, err := q.Exec(ctx, query, section.Title, section.Description, section.Position, status, section.ID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}

func (r *repositorySection) UpdateSectionStatus(ctx context.Context, sectionID int64, status string) error {
	const op = "repository.section.UpdateSectionStatus"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `UPDATE sections SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`

	tag, err := q.Exec(ctx, query, status, sectionID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrSectionNotFound)
	}
	return nil
}

func (r *repositorySection) UpdateStatusByCourseID(ctx context.Context, courseID int64, status string) error {
	const op = "repository.section.UpdateStatusByCourseID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `UPDATE sections SET status = $1, updated_at = NOW() WHERE course_id = $2 AND deleted_at IS NULL`

	_, err := q.Exec(ctx, query, status, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	return nil
}

func (r *repositorySection) ReorderSections(ctx context.Context, courseID int64, sectionIDs []int64) error {
	const op = "repository.section.ReorderSections"
	if len(sectionIDs) == 0 {
		return nil
	}

	positions := make([]int32, len(sectionIDs))
	for i := range sectionIDs {
		positions[i] = int32(i + 1)
	}

	query := `
		UPDATE sections AS s
		SET position = v.new_pos, updated_at = NOW()
		FROM (SELECT unnest($1::bigint[]) AS id, unnest($2::int[]) AS new_pos) AS v
		WHERE s.id = v.id AND s.course_id = $3 AND s.deleted_at IS NULL
	`

	q := txmanager.GetQueryExecutor(ctx, r.Pool)
	_, err := q.Exec(ctx, query, sectionIDs, positions, courseID)
	if err != nil {
		return fmt.Errorf("%s: batch update sections: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return nil
}

func (r *repositorySection) DeleteSection(ctx context.Context, sectionID int64) error {
	const op = "repository.section.DeleteSection"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

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

func (r *repositorySection) GetMaxPositionByCourseID(ctx context.Context, courseID int64) (int, error) {
	const op = "repository.section.GetMaxPositionByCourseID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `SELECT COALESCE(MAX(position), 0) FROM sections WHERE course_id = $1 AND deleted_at IS NULL`

	var maxPos int
	err := q.QueryRow(ctx, query, courseID).Scan(&maxPos)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return maxPos, nil
}
