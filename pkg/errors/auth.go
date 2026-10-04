package apperrors

import "errors"

var (
	ErrCreatedUser           = errors.New("failed to create user")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")
	ErrUserBanned            = errors.New("user is banned")
	ErrUserDeactivated       = errors.New("user deactivated")
	ErrRefreshTokenNotFound  = errors.New("refresh token not found")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrInvalidRole           = errors.New("invalid role")
	ErrCannotModifySelf      = errors.New("cannot modify own role or status")
)
