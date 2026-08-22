package response

import (
	errorsAPP "edtech/pkg/errors"
	"errors"
	"log/slog"
	"net/http"
)

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
		Error(w, r, http.StatusBadRequest, ValidationError, "Validation failed")
	// 400 Bad Request - ошибки URL параметров
	case errors.Is(err, errorsAPP.ErrInvalidURLParam):
		Error(w, r, http.StatusBadRequest, "INVALID_URL_PARAM", "Invalid URL parametr")
	case errors.Is(err, errorsAPP.ErrInvalidURLQuery):
		Error(w, r, http.StatusBadRequest, "INVALID_URL_QUERY", "Invalid URL query")

	// 401 Unauthorized - ошибки аутентификации
	case errors.Is(err, errorsAPP.ErrInvalidCredentials):
		Error(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	case errors.Is(err, errorsAPP.ErrInvalidJWT):
		Error(w, r, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
	case errors.Is(err, errorsAPP.ErrUnauthorized):
		Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
	case errors.Is(err, errorsAPP.ErrUserNotFound):
		Error(w, r, http.StatusNotFound, "USER_NOT_FOUND", "User not found")

	// 403 Forbidden - ошибки доступа
	case errors.Is(err, errorsAPP.ErrForbidden):
		Error(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied")
	case errors.Is(err, errorsAPP.ErrInvalidAction):
		Error(w, r, http.StatusForbidden, "INVALID_ACTION", "Action not allowed")

	// 404 Not Found
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
	case errors.Is(err, errorsAPP.ErrEnrolled):
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
