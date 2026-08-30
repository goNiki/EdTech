package quiz

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

type QuizHandler struct {
	quizService    service.QuizServices
	validator      validator.Validator
	authMiddleware auth.AuthMiddleware
}

func NewQuizHandler(quizService service.QuizServices, authMiddleware ...auth.AuthMiddleware) *QuizHandler {
	h := &QuizHandler{
		quizService: quizService,
		validator:   *validator.NewValidator(),
	}
	if len(authMiddleware) > 0 {
		h.authMiddleware = authMiddleware[0]
	}
	return h
}

func (h *QuizHandler) getUserID(ctx context.Context) int64 {
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
