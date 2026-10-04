package converter

import (
	"encoding/json"

	"edtech/internal/domain"
	repomodels "edtech/internal/repository/models"
)

func CreateUserToEntity(user domain.CreateUser) repomodels.User {
	return repomodels.User{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Username:     user.Username,
		Role:         string(user.Role),
	}
}

func UserToModel(user *repomodels.User) *domain.User {
	if user == nil {
		return nil
	}
	prefs := domain.DefaultUserPreferences()
	if len(user.Preferences) > 0 {
		_ = json.Unmarshal(user.Preferences, &prefs)
	}
	return &domain.User{
		ID:            user.ID,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		AvatarURL:     user.AvatarURL,
		Bio:           user.Bio,
		Headline:      user.Headline,
		Preferences:   prefs,
		Role:          domain.Role(user.Role),
		EmailVerified: user.EmailVerified,
		IsActive:      user.IsActive,
		IsBanned:      user.IsBanned,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		DeletedAt:     user.DeletedAt,
	}
}

func RefresTokenToDomain(token repomodels.RefreshTokenData) domain.RefreshTokenData {
	return domain.RefreshTokenData{
		ID:           token.ID,
		RefreshToken: token.RefreshToken,
		UserID:       token.UserID,
		ExpiresAt:    token.ExpiresAt,
		CreatedAt:    token.CreatedAt,
	}
}
