package category

import (
	"fmt"
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/validator"
	"edtech/internal/interfaces/handlers/converter"
	"edtech/internal/interfaces/middleware/auth"
	response "edtech/internal/interfaces/response"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

type CategoryHandler struct {
	categoryService service.CategoryServices
	validator       validator.Validator
	authMiddleware  auth.AuthMiddleware
}

func NewCategoryHandler(categoryService service.CategoryServices, authMiddleware auth.AuthMiddleware) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validator:       *validator.NewValidator(),
		authMiddleware:  authMiddleware,
	}
}

// ListCategories handles GET /api/v1/categories
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.category.ListCategories"
	log := logger.GetLogger(r.Context(), op)

	categories, err := h.categoryService.ListCategories(r.Context())
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, dto.CategoriesResponse{
		Categories: converter.CategoriesToDTO(categories),
	})
}

// CreateCategory handles POST /api/v1/categories (Admin/Auth)
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.category.CreateCategory"
	log := logger.GetLogger(r.Context(), op)

	var req dto.CreateCategoryRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	catDomain := converter.CreateCategoryRequestToDomain(req)
	id, err := h.categoryService.CreateCategory(r.Context(), &catDomain)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	catDomain.ID = id
	response.Created(w, r, converter.CategoryToDTO(&catDomain))
}
