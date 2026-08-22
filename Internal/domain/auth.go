package domain

import (
	errorsAPP "edtech/pkg/errors"
	"time"
)

type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleStudent, RoleTeacher, RoleAdmin:
		return true
	default:
		return false
	}
}


type User struct {
	ID            int64
	Email         string
	PasswordHash  string
	Username      string
	FirstName     *string
	LastName      *string
	AvatarURL     *string
	Bio           *string
	Role          Role
	EmailVerified bool
	IsActive      bool
	IsBanned      bool
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type CreateUser struct {
	Email        string
	PasswordHash string
	Username     string
	Role         Role
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) CanLogin() error {
	if u.IsBanned {
		return errorsAPP.ErrUserBanned
	}

	if !u.IsActive {
		return errorsAPP.ErrUserDeactivated
	}
	return nil
}

func (u *User) Login(now time.Time) {
	u.LastLoginAt = &now
}

type RegisterInput struct {
	Email    string
	Password string
	Username string
	Role     Role
}

type AuthTokens struct {
	ID           int64
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type RefreshTokenData struct {
	ID           int64
	RefreshToken string
	UserID       int64
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func (t *RefreshTokenData) IsValid(now time.Time) bool {
	return t.ExpiresAt.After(now)
}

type UpdateProfileInput struct {
	FirstName *string
	LastName  *string
	Bio       *string
	AvatarURL *string
}

func (u *User) UpdateProfile(input UpdateProfileInput) {
	if input.FirstName != nil {
		u.FirstName = input.FirstName
	}
	if input.LastName != nil {
		u.LastName = input.LastName
	}
	if input.Bio != nil {
		u.Bio = input.Bio
	}
	if input.AvatarURL != nil {
		u.AvatarURL = input.AvatarURL
	}
}

