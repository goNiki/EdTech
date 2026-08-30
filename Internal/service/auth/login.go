package auth

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"time"
)

func (s *service) Login(ctx context.Context, email, password string) (domain.AuthTokens, error) {
	const op = "service.auth.Login"

	user, err := s.repo.GetUserByEmail(ctx, s.db, email)
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := user.CanLogin(); err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	ok, err := s.hasherManager.CheckPassword(user.PasswordHash, password)

	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}
	if !ok {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidCredentials)
	}

	now := time.Now()
	tokens, err := s.generateAndSaveTokens(ctx, s.db, user, now)
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	user.Login(now)

	_ = s.repo.UpdateLastLogin(ctx, s.db, user.ID, now)

	return tokens, nil

}

func (s *service) generateAndSaveTokens(ctx context.Context, q db.QueryExecutor, user *domain.User, now time.Time) (domain.AuthTokens, error) {
	const op = "service.auth.generateAndSaveTokens"

	accessToken, err := s.jwtManager.GenerateAccessToken(user, now)
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, expiresIn, err := s.jwtManager.GenerateRefreshToken(user.ID, now)
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	hashRefreshToken := s.hasherManager.HashRefreshToken(refreshToken)

	if err := s.refreshRepo.DeleteAllByUserID(ctx, q, user.ID); err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.refreshRepo.Save(ctx, q, hashRefreshToken, user.ID, expiresIn); err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return domain.AuthTokens{
		ID:           user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtManager.GetAccessTokenExpiresIn(),
	}, nil

}
