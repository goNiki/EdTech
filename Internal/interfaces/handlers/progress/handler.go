package progress

import (
	"context"
	"net/http"
	"strconv"

	"edtech/internal/infrastructure/validator"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

type ProgressHandler struct {
	progressService service.ProgressServices
	validator       validator.Validator
	authMiddleware  auth.AuthMiddleware
}

func NewProgressHandler(progressService service.ProgressServices, authMiddleware ...auth.AuthMiddleware) *ProgressHandler {
	h := &ProgressHandler{
		progressService: progressService,
		validator:       *validator.NewValidator(),
	}
	if len(authMiddleware) > 0 {
		h.authMiddleware = authMiddleware[0]
	}
	return h
}

func (h *ProgressHandler) getUserID(ctx context.Context) int64 {
	if h.authMiddleware != nil {
		return h.authMiddleware.GetUserID(ctx)
	}
	return auth.GetUserID(ctx)
}

func parseIDParam(r *http.Request, keys ...string) (int64, error) {
	for _, key := range keys {
		val := chi.URLParam(r, key)
		if val != "" {
			id, err := strconv.ParseInt(val, 10, 64)
			if err != nil || id <= 0 {
				return 0, errorsAPP.ErrInvalidURLParam
			}
			return id, nil
		}
	}
	return 0, errorsAPP.ErrInvalidURLParam
}
