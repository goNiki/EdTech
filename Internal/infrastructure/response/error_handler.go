package response

import (
	errorsAPP "edtech/pkg/errors"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse структура ответа с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse структура успешного ответа
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// HandleError централизованная обработка ошибок
// Логирует ошибку с полным контекстом и возвращает подходящий HTTP статус
func HandleError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error, operation string) {
	if err == nil {
		return
	}

	// Всегда логируем с полным контекстом
	log.Error("operation failed",
		slog.String("operation", operation),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("error", err.Error()),
	)

	// Мапим внутренние ошибки на HTTP статусы
	switch {
	// 400 Bad Request - ошибки валидации
	case errors.Is(err, errorsAPP.ErrValidationFailed),
		errors.Is(err, errorsAPP.ErrCourseValidation),
		errors.Is(err, errorsAPP.ErrLessonValidation),
		errors.Is(err, errorsAPP.ErrFailValidate),
		errors.Is(err, errorsAPP.ErrFailEnroleValidate),
		errors.Is(err, errorsAPP.ErrNothingToUpdate):
		Error(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")

	// 401 Unauthorized - ошибки аутентификации
	case errors.Is(err, errorsAPP.ErrInvalidCredentials):
		Error(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	case errors.Is(err, errorsAPP.ErrInvalidJWT):
		Error(w, r, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")

	// 403 Forbidden - ошибки доступа
	case errors.Is(err, errorsAPP.ErrForbidden):
		Error(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied")
	case errors.Is(err, errorsAPP.ErrInvalidAction):
		Error(w, r, http.StatusForbidden, "INVALID_ACTION", "Action not allowed")

	// 404 Not Found
	case errors.Is(err, errorsAPP.ErrUserNotFound):
		Error(w, r, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
	case errors.Is(err, errorsAPP.ErrNotFoundCourse):
		Error(w, r, http.StatusNotFound, "COURSE_NOT_FOUND", "Course not found")
	case errors.Is(err, errorsAPP.ErrNotFoundLesson):
		Error(w, r, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson not found")

	// 409 Conflict - конфликты данных
	case errors.Is(err, errorsAPP.ErrSlugAlreadyExists):
		Error(w, r, http.StatusConflict, "SLUG_EXISTS", "Course with this slug already exists")
	case errors.Is(err, errorsAPP.ErrUserAlreadyEnrolled):
		Error(w, r, http.StatusConflict, "ALREADY_ENROLLED", "User already enrolled in this course")
	case errors.Is(err, errorsAPP.ErrCourseAlredyPublished):
		Error(w, r, http.StatusConflict, "ALREADY_PUBLISHED", "Course is already published")
	case errors.Is(err, errorsAPP.ErrEnrolledByCreated):
		Error(w, r, http.StatusConflict, "ENROLLMENT_ERROR", "Failed to enroll user")

	// 500 Internal Server Error - технические ошибки
	case errors.Is(err, errorsAPP.ErrInternalDB):
		Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	case errors.Is(err, errorsAPP.ErrFailHashingPassword):
		Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	case errors.Is(err, errorsAPP.ErrFailedCreateJWT),
		errors.Is(err, errorsAPP.ErrFailedCreateRefreshJwt):
		Error(w, r, http.StatusInternalServerError, "TOKEN_ERROR", "Failed to generate token")

	// По умолчанию - 500 Internal Server Error
	default:
		// НЕ раскрываем детали неизвестных ошибок
		Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}

// Error отправляет ответ с ошибкой
func Error(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	render.Status(r, status)
	render.JSON(w, r, ErrorResponse{
		Error:   message,
		Code:    code,
		Message: message,
	})
}

// ErrorWithDetails отправляет ответ с ошибкой и дополнительными деталями
func ErrorWithDetails(w http.ResponseWriter, r *http.Request, status int, code, message, details string) {
	render.Status(r, status)
	render.JSON(w, r, ErrorResponse{
		Error:   message,
		Code:    code,
		Message: details,
	})
}

// Success отправляет успешный ответ
func Success(w http.ResponseWriter, r *http.Request, status int, data interface{}, message string) {
	render.Status(r, status)
	render.JSON(w, r, SuccessResponse{
		Data:    data,
		Message: message,
	})
}

// OK отправляет успешный ответ с кодом 200
func OK(w http.ResponseWriter, r *http.Request, data interface{}) {
	Success(w, r, http.StatusOK, data, "")
}

// Created отправляет успешный ответ с кодом 201
func Created(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	Success(w, r, http.StatusCreated, data, message)
}

// NoContent отправляет успешный ответ без тела с кодом 204
func NoContent(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNoContent)
}
