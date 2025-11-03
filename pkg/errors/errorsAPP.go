package errorsAPP

import "errors"

// Domain errors (бизнес-логика)
var (
	// User errors
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyEnrolled = errors.New("user already enrolled")

	// Course errors
	ErrNotFoundCourse        = errors.New("course not found")
	ErrSlugAlreadyExists     = errors.New("slug already exists")
	ErrCourseAlredyPublished = errors.New("course already published")
	ErrCourseValidation      = errors.New("course validation failed")

	// Lesson errors
	ErrNotFoundLesson   = errors.New("lesson not found")
	ErrLessonValidation = errors.New("lesson validation failed")
	ErrNothingToUpdate  = errors.New("nothing to update")

	// Enrollment errors
	ErrEnrolledByCreated  = errors.New("failed to enroll creator")
	ErrFailEnroleValidate = errors.New("enrollment validation failed")

	// Permission errors
	ErrForbidden     = errors.New("access denied")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidAction = errors.New("invalid action")

	// Validation errors
	ErrValidationFailed = errors.New("validation failed")
	ErrFailValidate     = errors.New("failed to validate fields")
)

// Infrastructure errors (технические)
var (
	ErrInternalDB             = errors.New("internal database error")
	ErrFailHashingPassword    = errors.New("failed to hash password")
	ErrFailedCreateJWT        = errors.New("failed to create jwt token")
	ErrFailedCreateRefreshJwt = errors.New("failed to create refresh token")
	ErrInvalidJWT             = errors.New("invalid jwt token")
)
