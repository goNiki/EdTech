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
		errors.Is(err, errorsAPP.ErrQuizValidation),
		errors.Is(err, errorsAPP.ErrEmptyTitle),
		errors.Is(err, errorsAPP.ErrTitleTooLong),
		errors.Is(err, errorsAPP.ErrEmptySlug),
		errors.Is(err, errorsAPP.ErrInvalidSlug),
		errors.Is(err, errorsAPP.ErrFailValidate),
		errors.Is(err, errorsAPP.ErrFailEnroleValidate):
		Error(w, r, http.StatusBadRequest, ValidationError, "Validation failed")
	// 400 Bad Request - ошибки загрузки файлов
	case errors.Is(err, errorsAPP.ErrFileTooLarge):
		Error(w, r, http.StatusBadRequest, "FILE_TOO_LARGE", "File size exceeds maximum allowed limit (25MB)")
	case errors.Is(err, errorsAPP.ErrInvalidFileType):
		Error(w, r, http.StatusBadRequest, "INVALID_FILE_TYPE", "Unsupported or forbidden file type")
	case errors.Is(err, errorsAPP.ErrEmptyFile):
		Error(w, r, http.StatusBadRequest, "EMPTY_FILE", "Uploaded file cannot be empty")
	case errors.Is(err, errorsAPP.ErrCreatorCannotUnenroll):
		Error(w, r, http.StatusBadRequest, "CREATOR_CANNOT_UNENROLL", "Course creator cannot be unenrolled")
	case errors.Is(err, errorsAPP.ErrSamePassword):
		Error(w, r, http.StatusBadRequest, "SAME_PASSWORD", "New password cannot be the same as old password")
	case errors.Is(err, errorsAPP.ErrPasswordTooShort):
		Error(w, r, http.StatusBadRequest, "PASSWORD_TOO_SHORT", "New password must be at least 8 characters")
	// 400 Bad Request - ошибки URL параметров
	case errors.Is(err, errorsAPP.ErrInvalidURLParam):
		Error(w, r, http.StatusBadRequest, "INVALID_URL_PARAM", "Invalid URL parametr")
	case errors.Is(err, errorsAPP.ErrInvalidURLQuery):
		Error(w, r, http.StatusBadRequest, "INVALID_URL_QUERY", "Invalid URL query")
	case errors.Is(err, errorsAPP.ErrTimeLimitExceeded):
		Error(w, r, http.StatusBadRequest, "TIME_LIMIT_EXCEEDED", "Quiz time limit exceeded")

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
	case errors.Is(err, errorsAPP.ErrNotEnrolled):
		Error(w, r, http.StatusForbidden, "NOT_ENROLLED", "User not enrolled in course")

	// 404 Not Found
	case errors.Is(err, errorsAPP.ErrNotFoundCourse), errors.Is(err, errorsAPP.ErrCourseNotFound):
		Error(w, r, http.StatusNotFound, "COURSE_NOT_FOUND", "Course not found")
	case errors.Is(err, errorsAPP.ErrCategoryNotFound):
		Error(w, r, http.StatusNotFound, "CATEGORY_NOT_FOUND", "Category not found")
	case errors.Is(err, errorsAPP.ErrNotFoundLesson):
		Error(w, r, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson not found")
	case errors.Is(err, errorsAPP.ErrProgressNotFound),
		errors.Is(err, errorsAPP.ErrCourseProgressNotFound),
		errors.Is(err, errorsAPP.ErrLessonProgressNotFound):
		Error(w, r, http.StatusNotFound, "PROGRESS_NOT_FOUND", "Progress not found")
	case errors.Is(err, errorsAPP.ErrQuizNotFound), errors.Is(err, errorsAPP.ErrNotFoundQuiz):
		Error(w, r, http.StatusNotFound, "QUIZ_NOT_FOUND", "Quiz not found")
	case errors.Is(err, errorsAPP.ErrAttemptNotFound), errors.Is(err, errorsAPP.ErrNotFoundAttempt):
		Error(w, r, http.StatusNotFound, "ATTEMPT_NOT_FOUND", "Quiz attempt not found")
	case errors.Is(err, errorsAPP.ErrAnswerNotFound):
		Error(w, r, http.StatusNotFound, "ANSWER_NOT_FOUND", "Quiz answer not found")

	// 409 Conflict - конфликты данных
	case errors.Is(err, errorsAPP.ErrNothingToUpdate):
		Error(w, r, http.StatusConflict, "CONFLICT", "Resource was updated by another request or no changes were made")
	case errors.Is(err, errorsAPP.ErrEmailAlreadyExists):
		Error(w, r, http.StatusConflict, "EMAIL_EXISTS", "User with this email already exists")
	case errors.Is(err, errorsAPP.ErrUserNameAlreadyExists):
		Error(w, r, http.StatusConflict, "USERNAME_EXISTS", "User with this username already exists")
	case errors.Is(err, errorsAPP.ErrSlugAlreadyExists):
		Error(w, r, http.StatusConflict, "SLUG_EXISTS", "Course with this slug already exists")
	case errors.Is(err, errorsAPP.ErrUserAlreadyEnrolled):
		Error(w, r, http.StatusConflict, "ALREADY_ENROLLED", "User already enrolled in this course")
	case errors.Is(err, errorsAPP.ErrCourseAlreadyPublished):
		Error(w, r, http.StatusConflict, "ALREADY_PUBLISHED", "Course is already published")
	case errors.Is(err, errorsAPP.ErrCourseAlreadyArchived):
		Error(w, r, http.StatusConflict, "ALREADY_ARCHIVED", "Course is already archived")
	case errors.Is(err, errorsAPP.ErrCannotPublishEmptyCourse):
		Error(w, r, http.StatusConflict, "CANNOT_PUBLISH_EMPTY_COURSE", "Cannot publish course without lessons")
	case errors.Is(err, errorsAPP.ErrEnrolled):
		Error(w, r, http.StatusConflict, "ENROLLMENT_ERROR", "Failed to enroll user")
	case errors.Is(err, errorsAPP.ErrMaxAttemptsReached):
		Error(w, r, http.StatusConflict, "MAX_ATTEMPTS_REACHED", "Maximum quiz attempts reached")
	case errors.Is(err, errorsAPP.ErrAttemptAlreadyCompleted):
		Error(w, r, http.StatusConflict, "ATTEMPT_ALREADY_COMPLETED", "Quiz attempt already completed")
	case errors.Is(err, errorsAPP.ErrCannotEnrollStudentInDraft):
		Error(w, r, http.StatusConflict, "CANNOT_ENROLL_STUDENT_IN_DRAFT", "Cannot enroll student in a draft course")

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
