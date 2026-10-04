package category

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repositoryImpl) ListCategories(ctx context.Context, q db.QueryExecutor) ([]domain.Category, error) {
	const op = "repository.category.ListCategories"

	query := `
		SELECT 
			c.id, 
			c.name, 
			c.slug, 
			COALESCE(c.description, ''), 
			COALESCE(c.icon_url, ''),
			c.parent_id,
			c.position,
			COUNT(crs.id) FILTER (WHERE crs.status = 'published' AND crs.deleted_at IS NULL) AS courses_count
		FROM categories c
		LEFT JOIN courses crs ON crs.category_id = c.id
		GROUP BY c.id, c.name, c.slug, c.description, c.icon_url, c.parent_id, c.position
		ORDER BY c.position ASC, c.id ASC`

	rows, err := q.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Slug,
			&cat.Description,
			&cat.IconURL,
			&cat.ParentID,
			&cat.Position,
			&cat.CoursesCount,
		); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if categories == nil {
		categories = []domain.Category{}
	}

	return categories, nil
}

func (r *repositoryImpl) CreateCategory(ctx context.Context, q db.QueryExecutor, category *domain.Category) (*domain.Category, error) {
	const op = "repository.category.CreateCategory"

	query := `
		INSERT INTO categories (name, slug, description, icon_url, parent_id, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id`

	err := q.QueryRow(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.IconURL,
		category.ParentID,
		category.Position,
	).Scan(&category.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return category, nil
}

func (r *repositoryImpl) GetCategoryByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Category, error) {
	const op = "repository.category.GetCategoryByID"

	query := `
		SELECT id, name, slug, COALESCE(description, ''), COALESCE(icon_url, ''), parent_id, position
		FROM categories
		WHERE id = $1`

	var cat domain.Category
	err := q.QueryRow(ctx, query, id).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&cat.Description,
		&cat.IconURL,
		&cat.ParentID,
		&cat.Position,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCategoryNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return &cat, nil
}

func (r *repositoryImpl) GetCategoryBySlug(ctx context.Context, q db.QueryExecutor, slug string) (*domain.Category, error) {
	const op = "repository.category.GetCategoryBySlug"

	query := `
		SELECT id, name, slug, COALESCE(description, ''), COALESCE(icon_url, ''), parent_id, position
		FROM categories
		WHERE slug = $1`

	var cat domain.Category
	err := q.QueryRow(ctx, query, slug).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&cat.Description,
		&cat.IconURL,
		&cat.ParentID,
		&cat.Position,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCategoryNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return &cat, nil
}

func (r *repositoryImpl) UpdateCategory(ctx context.Context, q db.QueryExecutor, category *domain.Category) error {
	const op = "repository.category.UpdateCategory"

	query := `
		UPDATE categories
		SET name = $2, slug = $3, description = $4, icon_url = $5, parent_id = $6, position = $7, updated_at = NOW()
		WHERE id = $1`

	tag, err := q.Exec(ctx, query,
		category.ID,
		category.Name,
		category.Slug,
		category.Description,
		category.IconURL,
		category.ParentID,
		category.Position,
	)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrCategoryNotFound)
	}

	return nil
}

func (r *repositoryImpl) DeleteCategory(ctx context.Context, q db.QueryExecutor, id int64) error {
	const op = "repository.category.DeleteCategory"

	query := `DELETE FROM categories WHERE id = $1`

	tag, err := q.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrCategoryNotFound)
	}

	return nil
}
