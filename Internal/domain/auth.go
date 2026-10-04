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
	Headline      *string
	Preferences   UserPreferences
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
	Headline  *string
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
	if input.Headline != nil {
		u.Headline = input.Headline
	}
}

type UserPreferences struct {
	FontScale    string `json:"font_scale"`
	ContentWidth string `json:"content_width"`
	LineHeight   string `json:"line_height"`
	ReadingTheme string `json:"reading_theme"`
}

func DefaultUserPreferences() UserPreferences {
	return UserPreferences{
		FontScale:    "medium",
		ContentWidth: "wide",
		LineHeight:   "normal",
		ReadingTheme: "system",
	}
}

type UpdatePreferencesInput struct {
	FontScale    *string
	ContentWidth *string
	LineHeight   *string
	ReadingTheme *string
}

func (i UpdatePreferencesInput) Validate() error {
	if i.FontScale != nil {
		switch *i.FontScale {
		case "compact", "medium", "large", "xlarge":
		default:
			return errorsAPP.ErrInvalidPreferences
		}
	}
	if i.ContentWidth != nil {
		switch *i.ContentWidth {
		case "standard", "wide", "full":
		default:
			return errorsAPP.ErrInvalidPreferences
		}
	}
	if i.LineHeight != nil {
		switch *i.LineHeight {
		case "normal", "relaxed":
		default:
			return errorsAPP.ErrInvalidPreferences
		}
	}
	if i.ReadingTheme != nil {
		switch *i.ReadingTheme {
		case "system", "light", "dark", "sepia":
		default:
			return errorsAPP.ErrInvalidPreferences
		}
	}
	return nil
}

func (u *User) UpdatePreferences(input UpdatePreferencesInput) {
	if input.FontScale != nil {
		u.Preferences.FontScale = *input.FontScale
	}
	if input.ContentWidth != nil {
		u.Preferences.ContentWidth = *input.ContentWidth
	}
	if input.LineHeight != nil {
		u.Preferences.LineHeight = *input.LineHeight
	}
	if input.ReadingTheme != nil {
		u.Preferences.ReadingTheme = *input.ReadingTheme
	}
}

type UserFilter struct {
	Search   *string
	Role     *Role
	IsBanned *bool
}

