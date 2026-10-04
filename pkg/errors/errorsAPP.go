package errorsAPP

import "errors"

var (
	// user Errors
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyEnrolled = errors.New("user already enrolled")
	ErrSamePassword        = errors.New("new password cannot be the same as old password")
	ErrPasswordTooShort    = errors.New("new password must be at least 8 characters")

	// General / review / certificate errors
	ErrNotFound            = errors.New("not found")
	ErrReviewNotFound      = errors.New("review not found")
	ErrCertificateNotFound = errors.New("certificate not found")
	ErrCourseNotCompleted  = errors.New("course not completed")

	// Course errors
	ErrCourseNotFound           = errors.New("course not found")
	ErrNotFoundCourse           = ErrCourseNotFound
	ErrSlugAlreadyExists        = errors.New("slug already exists")
	ErrCourseAlreadyPublished    = errors.New("course already published")
	ErrCourseAlreadyArchived     = errors.New("course already archived")
	ErrCannotPublishEmptyCourse = errors.New("cannot publish course without lessons")
	ErrCourseValidation         = errors.New("course validation failed")
	ErrEmptyTitle               = errors.New("title cannot be empty")
	ErrTitleTooLong             = errors.New("title is too long")
	ErrEmptySlug                = errors.New("slug cannot be empty")
	ErrInvalidSlug              = errors.New("invalid slug format")

	// Category errors
	ErrCategoryNotFound         = errors.New("category not found")

	// Section errors
	ErrSectionNotFound          = errors.New("section not found")
	ErrSectionAlreadyPublished  = errors.New("section already published")
	ErrSectionAlreadyArchived   = errors.New("section already archived")

	// Lesson errors
	ErrLessonNotFound           = errors.New("lesson not found")
	ErrNotFoundLesson           = ErrLessonNotFound
	ErrLessonAlreadyPublished   = errors.New("lesson already published")
	ErrLessonAlreadyArchived    = errors.New("lesson already archived")
	ErrLessonValidation         = errors.New("lesson validation failed")
	ErrNothingToUpdate          = errors.New("nothing to update")

	// Progress errors
	ErrProgressNotFound       = errors.New("progress not found")
	ErrCourseProgressNotFound = errors.New("course progress not found")
	ErrLessonProgressNotFound = errors.New("lesson progress not found")

	// Quiz errors
	ErrQuizNotFound            = errors.New("quiz not found")
	ErrNotFoundQuiz            = ErrQuizNotFound
	ErrQuizValidation          = errors.New("quiz validation failed")
	ErrAttemptNotFound         = errors.New("quiz attempt not found")
	ErrNotFoundAttempt         = ErrAttemptNotFound
	ErrMaxAttemptsReached      = errors.New("maximum quiz attempts reached")
	ErrTimeLimitExceeded       = errors.New("quiz time limit exceeded")
	ErrAttemptAlreadyCompleted = errors.New("quiz attempt already completed")
	ErrAnswerNotFound          = errors.New("quiz answer not found")

	// Enrollment errors
	ErrEnrolled                   = errors.New("failed to enroll")
	ErrFailEnroleValidate         = errors.New("enrollment validation failed")
	ErrNotEnrolled                = errors.New("user not enrolled in course")
	ErrCreatorCannotUnenroll      = errors.New("creator cannot unenroll from course")
	ErrCannotEnrollStudentInDraft = errors.New("cannot enroll student in a draft course")

	// Permission errors
	ErrForbidden           = errors.New("access denied")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidAction       = errors.New("invalid action")
	ErrCheckingPermissions = errors.New("error checking permissions")

	// Upload errors
	ErrFileTooLarge     = errors.New("file size exceeds maximum allowed limit")
	ErrInvalidFileType  = errors.New("unsupported or dangerous file type")
	ErrEmptyFile        = errors.New("uploaded file cannot be empty")
	ErrUploadFailed     = errors.New("failed to upload file")
	ErrPathTraversal    = errors.New("security violation: path traversal detected")

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
