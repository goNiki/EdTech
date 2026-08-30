package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func UserToDTO(d *domain.User) dto.User {
	if d == nil {
		return dto.User{}
	}
	return dto.User{
		ID:            d.ID,
		Email:         d.Email,
		Username:      d.Username,
		FirstName:     d.FirstName,
		LastName:      d.LastName,
		AvatarURL:     d.AvatarURL,
		Bio:           d.Bio,
		Role:          string(d.Role),
		EmailVerified: d.EmailVerified,
		IsActive:      d.IsActive,
		LastLoginAt:   d.LastLoginAt,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

func RegisterRequestToDomain(req dto.RegisterRequest) domain.User {
	return domain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		Username:     req.Username,
		Role:         domain.RoleStudent,
	}
}

func RegisterRequestToInput(req dto.RegisterRequest) domain.RegisterInput {
	return domain.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	}
}

func UpdateProfileRequestToDomain(req dto.UpdateProfileRequest) domain.UpdateProfileInput {
	return domain.UpdateProfileInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
	}
}
