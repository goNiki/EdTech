package dto

import (
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"time"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type RegisterResponse struct {
	User User `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=1,max=100"`
}

type LoginResponse struct {
	ID           int64  `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type EnrollRequest struct {
	UserEmail string `json:"user_email" validate:"required,email,max=255"`
	CourseID  int64  `json:"course_id" validate:"required,gt=0"`
	Role      string `json:"role" validate:"required,oneof=student teacher creator"`
}

type User struct {
	ID            int64      `json:"id"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	FirstName     *string    `json:"first_name,omitempty"`
	LastName      *string    `json:"last_name,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	Bio           *string    `json:"bio,omitempty"`
	Headline      *string            `json:"headline,omitempty"`
	Preferences   UserPreferencesDTO `json:"preferences"`
	Role          string             `json:"role"`
	EmailVerified bool       `json:"email_verified"`
	IsActive      bool       `json:"is_active"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type UserRole struct {
	UserID   int64
	UserRole string
}

func (u *UserRole) Validate() error {
	if u.UserID == 0 {
		return errorsAPP.ErrUnauthorized
	}
	if u.UserRole == "" {
		return errorsAPP.ErrForbidden
	}
	return nil
}

type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	Bio       *string `json:"bio,omitempty" validate:"omitempty,max=1000"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,http_url,max=500"`
	Headline  *string `json:"headline,omitempty" validate:"omitempty,max=150"`
}

type UpdateProfileResponse struct {
	User User `json:"user"`
}

type UserPreferencesDTO struct {
	FontScale    string `json:"font_scale"`
	ContentWidth string `json:"content_width"`
	LineHeight   string `json:"line_height"`
	ReadingTheme string `json:"reading_theme"`
}

type UpdatePreferencesRequest struct {
	FontScale    *string `json:"font_scale,omitempty" validate:"omitempty,oneof=compact medium large xlarge"`
	ContentWidth *string `json:"content_width,omitempty" validate:"omitempty,oneof=standard wide full"`
	LineHeight   *string `json:"line_height,omitempty" validate:"omitempty,oneof=normal relaxed"`
	ReadingTheme *string `json:"reading_theme,omitempty" validate:"omitempty,oneof=system light dark sepia"`
}

type UpdatePreferencesData struct {
	Preferences UserPreferencesDTO `json:"preferences"`
}

type UpdatePreferencesResponse struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    UpdatePreferencesData `json:"data"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=10,max=512"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=10,max=512"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=1,max=100"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100,nefield=OldPassword"`
}

type ChangeUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=student teacher admin"`
}

type SetUserBannedRequest struct {
	IsBanned bool `json:"is_banned"`
}

type AdminUserItemResponse struct {
	ID        int64       `json:"id"`
	Email     string      `json:"email"`
	Username  string      `json:"username"`
	FirstName *string     `json:"first_name,omitempty"`
	LastName  *string     `json:"last_name,omitempty"`
	Role      domain.Role `json:"role"`
	IsBanned  bool        `json:"is_banned"`
	CreatedAt string      `json:"created_at"`
}

type AdminUsersListResponse struct {
	Users    []AdminUserItemResponse `json:"users"`
	Total    int64                   `json:"total"`
	Page     int64                   `json:"page"`
	PageSize int64                   `json:"page_size"`
}
