package errorsAPP

import "errors"

var (
	ErrFailValidate           = errors.New("failed validate fields")
	ErrFailHashingPassword    = errors.New("fieled hash password")
	ErrUserNotFound           = errors.New("users not found")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrFailedCreateJWT        = errors.New("failed to create jwt token")
	ErrFailedCreateRefreshJwt = errors.New("failed to create refresh token")
	ErrInvalidJWT             = errors.New(("invalid jwt token"))
	ErrCourseValidation       = errors.New("course validation failed")
	ErrSlugAlreadyExists      = errors.New("slug already exists")
	ErrNotFoundCourse         = errors.New("course is not found ")
	ErrInternalDB             = errors.New("internal database error")
	ErrCourseAlredyPublished  = errors.New("course already published")
	ErrLessonValidation       = errors.New("lesson validation failed")
	ErrNotFoundLesson         = errors.New("lesson is not found")
	ErrNothingToUpdate        = errors.New("nothing to update for lesson")
	ErrEnrolledByCreated      = errors.New("error enrolled by created")
	ErrFailEnroleValidate     = errors.New("enrole validation failed")
	ErrUserAlreadyEnrolled    = errors.New("user already enrolled")
)
