package category_test

import (
	"context"
	"errors"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	categoryService "edtech/internal/service/category"
	errorsAPP "edtech/pkg/errors"
)

type mockCategoryRepo struct {
	repository.CategoryRepository
	categories []domain.Category
	listErr    error
	createdID  int64
	createErr  error
	category   *domain.Category
	getErr     error
}

func (m *mockCategoryRepo) ListCategories(ctx context.Context) ([]domain.Category, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.categories, nil
}

func (m *mockCategoryRepo) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	category.ID = m.createdID
	return category, nil
}

func (m *mockCategoryRepo) GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.category, nil
}

func TestListCategories_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockCategoryRepo{
		categories: []domain.Category{
			{
				ID:           1,
				Name:         "Программирование",
				Slug:         "programming",
				Description:  "Курсы по разработке",
				IconURL:      "https://example.com/icon.svg",
				CoursesCount: 12,
			},
			{
				ID:           2,
				Name:         "Дизайн",
				Slug:         "design",
				Description:  "UI/UX дизайн",
				CoursesCount: 5,
			},
		},
	}

	svc := categoryService.NewCategoryService(mockRepo)

	cats, err := svc.ListCategories(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cats) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(cats))
	}

	if cats[0].Name != "Программирование" || cats[0].CoursesCount != 12 {
		t.Errorf("unexpected category 0: %+v", cats[0])
	}
	if cats[1].Name != "Дизайн" || cats[1].CoursesCount != 5 {
		t.Errorf("unexpected category 1: %+v", cats[1])
	}
}

func TestListCategories_RepoError(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockCategoryRepo{
		listErr: errors.New("db error"),
	}

	svc := categoryService.NewCategoryService(mockRepo)

	_, err := svc.ListCategories(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateCategory_Validation(t *testing.T) {
	ctx := context.Background()
	svc := categoryService.NewCategoryService(&mockCategoryRepo{})

	_, err := svc.CreateCategory(ctx, &domain.Category{Name: ""})
	if !errors.Is(err, errorsAPP.ErrValidationFailed) {
		t.Errorf("expected ErrValidationFailed for empty name, got %v", err)
	}

	_, err = svc.CreateCategory(ctx, &domain.Category{Name: "Test", Slug: ""})
	if !errors.Is(err, errorsAPP.ErrValidationFailed) {
		t.Errorf("expected ErrValidationFailed for empty slug, got %v", err)
	}
}

func TestCreateCategory_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockCategoryRepo{
		createdID: 42,
	}

	svc := categoryService.NewCategoryService(mockRepo)

	id, err := svc.CreateCategory(ctx, &domain.Category{
		Name: "Data Science",
		Slug: "data-science",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != 42 {
		t.Errorf("expected id 42, got %d", id)
	}
}

func TestGetCategoryByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockCategoryRepo{
		category: &domain.Category{ID: 1, Name: "Frontend"},
	}

	svc := categoryService.NewCategoryService(mockRepo)

	cat, err := svc.GetCategoryByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.ID != 1 || cat.Name != "Frontend" {
		t.Errorf("unexpected category: %+v", cat)
	}
}
