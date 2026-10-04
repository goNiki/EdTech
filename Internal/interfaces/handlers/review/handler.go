package review

import (
	"log/slog"
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/interfaces/response"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ReviewHandler struct {
	reviewService  service.ReviewServices
	log            *slog.Logger
	validator      *validator.Validate
	authMiddleware auth.AuthMiddleware
}

func NewReviewHandler(
	reviewService service.ReviewServices,
	log *slog.Logger,
	validator *validator.Validate,
	authMiddleware auth.AuthMiddleware,
) *ReviewHandler {
	return &ReviewHandler{
		reviewService:  reviewService,
		log:            log,
		validator:      validator,
		authMiddleware: authMiddleware,
	}
}

// ListReviews handles GET /api/v1/courses/{courseid}/reviews
func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.review.ListReviews"

	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseid"), 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	page := 1
	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if psStr := r.URL.Query().Get("pagesize"); psStr != "" {
		if ps, err := strconv.Atoi(psStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	summary, err := h.reviewService.ListCourseReviews(r.Context(), courseID, page, pageSize)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	reviewDTOs := make([]dto.ReviewResponse, 0, len(summary.Reviews))
	for _, rev := range summary.Reviews {
		reviewDTOs = append(reviewDTOs, dto.ReviewResponse{
			ID:         rev.ID,
			CourseID:   rev.CourseID,
			UserID:     rev.UserID,
			Rating:     rev.Rating,
			Comment:    rev.Comment,
			CreatedAt:  rev.CreatedAt,
			UpdatedAt:  rev.UpdatedAt,
			UserName:   rev.UserName,
			UserAvatar: rev.UserAvatar,
		})
	}

	response.OK(w, r, dto.CourseReviewsListResponse{
		Reviews:       reviewDTOs,
		AverageRating: summary.AverageRating,
		ReviewsCount:  summary.ReviewsCount,
		Total:         summary.Total,
	})
}

// AddOrUpdateReview handles POST /api/v1/courses/{courseid}/reviews
func (h *ReviewHandler) AddOrUpdateReview(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.review.AddOrUpdateReview"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseid"), 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.CreateReviewRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrValidationFailed, op)
		return
	}

	rev, err := h.reviewService.AddOrUpdateReview(r.Context(), userID, courseID, req.Rating, req.Comment)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, dto.ReviewResponse{
		ID:        rev.ID,
		CourseID:  rev.CourseID,
		UserID:    rev.UserID,
		Rating:    rev.Rating,
		Comment:   rev.Comment,
		CreatedAt: rev.CreatedAt,
		UpdatedAt: rev.UpdatedAt,
	})
}

// DeleteReview handles DELETE /api/v1/courses/{courseid}/reviews
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.review.DeleteReview"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseid"), 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	if err := h.reviewService.DeleteReview(r.Context(), userID, courseID); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "review successfully deleted"})
}

// GetMyReview handles GET /api/v1/courses/{courseid}/reviews/my
func (h *ReviewHandler) GetMyReview(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.review.GetMyReview"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseid"), 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	rev, err := h.reviewService.GetMyReview(r.Context(), userID, courseID)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, dto.ReviewResponse{
		ID:        rev.ID,
		CourseID:  rev.CourseID,
		UserID:    rev.UserID,
		Rating:    rev.Rating,
		Comment:   rev.Comment,
		CreatedAt: rev.CreatedAt,
		UpdatedAt: rev.UpdatedAt,
	})
}
