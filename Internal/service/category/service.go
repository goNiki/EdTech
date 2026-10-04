package category

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository"
	services "edtech/internal/service"
	errorsAPP "edtech/pkg/errors"
)

var _ services.CategoryServices = (*serviceImpl)(nil)

type serviceImpl struct {
	repo repository.CategoryRepository
	db   db.QueryExecutor
}

func NewCategoryService(repo repository.CategoryRepository, database db.QueryExecutor) services.CategoryServices {
	return &serviceImpl{
		repo: repo,
		db:   database,
	}
}

func (s *serviceImpl) ListCategories(ctx context.Context) ([]domain.Category, error) {
	const op = "service.category.ListCategories"
	cats, err := s.repo.ListCategories(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return cats, nil
}

func (s *serviceImpl) CreateCategory(ctx context.Context, category *domain.Category) (int64, error) {
	const op = "service.category.CreateCategory"
	if category.Name == "" || category.Slug == "" {
		return 0, fmt.Errorf("%s: %w", op, errorsAPP.ErrValidationFailed)
	}
	created, err := s.repo.CreateCategory(ctx, s.db, category)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return created.ID, nil
}

func (s *serviceImpl) GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error) {
	const op = "service.category.GetCategoryByID"
	cat, err := s.repo.GetCategoryByID(ctx, s.db, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return cat, nil
}

func (s *serviceImpl) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	const op = "service.category.GetCategoryBySlug"
	cat, err := s.repo.GetCategoryBySlug(ctx, s.db, slug)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return cat, nil
}

func (s *serviceImpl) UpdateCategory(ctx context.Context, category *domain.Category) error {
	const op = "service.category.UpdateCategory"
	if err := s.repo.UpdateCategory(ctx, s.db, category); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *serviceImpl) DeleteCategory(ctx context.Context, id int64) error {
	const op = "service.category.DeleteCategory"
	if err := s.repo.DeleteCategory(ctx, s.db, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
