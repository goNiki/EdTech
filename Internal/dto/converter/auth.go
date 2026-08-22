package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func RegisterRequestoModels(req dto.RegisterRequest) domain.RegisterInput {

	return domain.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
		Role:     domain.RoleStudent,
	}

}

func UserToDTO(user domain.User) dto.User {
	return dto.User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
