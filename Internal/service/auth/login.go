package auth

import (
	"context"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
)

// TODO подумать над тем где и как правильно делать логирование и проверку ошибок
func (s *service) Login(ctx context.Context, email, password string) (*dto.LoginResponce, error) {
	const op = "usercase.auth.Login"

	log := logger.GetLogger(ctx, op)

	//TODO сделать обработку ошибок если пользователя нет, если ошибка с БД, и т.д.
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("Ошибка пользователя", sl.Error(err))
		return nil, err
	}

	if !s.hasherManager.CheckPassword(user.PasswordHash, password) {
		return nil, errorsAPP.ErrInvalidCredentials
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		log.Error("Error generate AccessToketn ", sl.Error(err))
		return nil, err
	}

	refreshToken, expiresAt, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Error("Error generate RefreshToken", sl.Error(err))
		return nil, err
	}

	hashRefreshToken, err := s.hasherManager.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	err = s.refreshRepo.Save(ctx, hashRefreshToken, user.ID, expiresAt)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponce{
		ID:           user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresID:    s.jwtManager.GetAccessTokenExpiresIn(),
	}, err

}
