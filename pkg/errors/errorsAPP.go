package errorsAPP

import "errors"

var (
	// user Errors
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyEnrolled = errors.New("user already enrolled")


	// Course errors
	ErrNotFoundCourse        = errors.New("course not found")
	ErrSlugAlreadyExists     = errors.New("slug already exists")
	ErrCourseAlredyPublished = errors.New("course already published")
	ErrCourseValidation      = errors.New("course validation failed")
	ErrEmptyTitle            = errors.New("title cannot be empty")
	ErrTitleTooLong          = errors.New("title is too long")
	ErrEmptySlug             = errors.New("slug cannot be empty")
	ErrInvalidSlug           = errors.New("invalid slug format")

	// Lesson errors
	ErrNotFoundLesson   = errors.New("lesson not found")
	ErrLessonValidation = errors.New("lesson validation failed")
	ErrNothingToUpdate  = errors.New("nothing to update")

	// Enrollment errors
	ErrEnrolled           = errors.New("failed to enroll")
	ErrFailEnroleValidate = errors.New("enrollment validation failed")

	// Permission errors
	ErrForbidden           = errors.New("access denied")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidAction       = errors.New("invalid action")
	ErrCheckingPermissions = errors.New("error checking permissions")

	// Validation errors
	ErrValidationFailed = errors.New("validation failed")
	ErrFailValidate     = errors.New("failed to validate fields")
)

// Infrastructure errors (технические)
var (
	ErrInternalDB             = errors.New("internal database error")
	ErrDecodeJSON             = errors.New("failed to decode json ")
	ErrFailHashingPassword    = errors.New("failed to hash password")
	ErrFailedCreateJWT        = errors.New("failed to create jwt token")
	ErrFailedCreateRefreshJwt = errors.New("failed to create refresh token")
	ErrInvalidJWT             = errors.New("invalid jwt token")
	ErrParseLoggerConfig      = errors.New("failed to parse logger config")
	ErrParsePostgresConfig    = errors.New("failed to parse postgres config")
	ErrParseServerConfig      = errors.New("failed to parse server config")
	ErrParseJWTConfig         = errors.New("failed to parse jwt config")
	ErrLoadEnv                = errors.New("failed to load env file")
	ErrFailCheckingPassword   = errors.New("failed to check password")
	ErrSaveRefreshToken       = errors.New("failed to save refresh token")
)

// Handlers errors
var (
	ErrInvalidURLParam = errors.New("invalid URL param")
	ErrInvalidURLQuery = errors.New("invalid URL query")
)

var (
	ErrSetGooseDialect = errors.New("failed to set goose dialect")
	ErrCloseDb         = errors.New("failed to close db connection")
	ErrUpMigration     = errors.New("failed to up migration")
	ErrUpToMigration   = errors.New("failed to upto migration")
	ErrDownMigration   = errors.New("failed to down migration")
	ErrDownToMigration = errors.New("failed to downto migration")
	ErrCreateMigration = errors.New("failed to create migration")
	ErrGetDbVersion    = errors.New("failed to get db version")
	ErrGetGooseStatus  = errors.New("failed to get goose status")
)
