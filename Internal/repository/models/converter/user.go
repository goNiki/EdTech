package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func UserToDomain(m *models.User) *domain.User {
	if m == nil {
		return nil
	}
	return &domain.User{
		ID:            m.ID,
		Email:         m.Email,
		PasswordHash:  m.PasswordHash,
		Username:      m.Username,
		FirstName:     m.FirstName,
		LastName:      m.LastName,
		AvatarURL:     m.AvatarURL,
		Bio:           m.Bio,
		Role:          domain.Role(m.Role),
		EmailVerified: m.EmailVerified,
		IsActive:      m.IsActive,
		IsBanned:      m.IsBanned,
		LastLoginAt:   m.LastLoginAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
}

func UserToEntity(d *domain.User) *models.User {
	if d == nil {
		return nil
	}
	return &models.User{
		ID:            d.ID,
		Email:         d.Email,
		PasswordHash:  d.PasswordHash,
		Username:      d.Username,
		FirstName:     d.FirstName,
		LastName:      d.LastName,
		AvatarURL:     d.AvatarURL,
		Bio:           d.Bio,
		Role:          string(d.Role),
		EmailVerified: d.EmailVerified,
		IsActive:      d.IsActive,
		IsBanned:      d.IsBanned,
		LastLoginAt:   d.LastLoginAt,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
		DeletedAt:     d.DeletedAt,
	}
}
