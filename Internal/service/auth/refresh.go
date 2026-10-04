package auth

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (domain.AuthTokens, error) {
	const op = "service.auth.RefreshToken"

	hashRefreshToken := s.hasherManager.HashRefreshToken(refreshToken)

	tokenData, err := s.refreshRepo.GetByToken(ctx, hashRefreshToken)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrRefreshTokenNotFound) {
			return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidRefreshToken)
		}
		return domain.AuthTokens{}, fmt.Errorf("%s: %w: %v", op, errorsAPP.ErrInternalDB, err)
	}

	now := time.Now()

	if !tokenData.IsValid(now) {
		if err := s.refreshRepo.DeleteRefreshToken(ctx, hashRefreshToken); err != nil {
			slog.Warn("failed to delete expired refresh token",
				"error", err,
				"op", op,
			)
		}
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidRefreshToken)
	}

	user, err := s.repo.GetUserByID(ctx, tokenData.UserID)
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := user.CanLogin(); err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	var newTokens domain.AuthTokens

	err = s.txmanager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {

		if err := s.refreshRepo.DeleteRefreshToken(ctx, hashRefreshToken); err != nil {
			return err
		}

		var err error
		newTokens, err = s.generateAndSaveTokens(ctx, user, now)
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return newTokens, nil
}
