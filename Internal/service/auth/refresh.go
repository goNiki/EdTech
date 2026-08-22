package auth

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (domain.AuthTokens, error) {
	const op = "service.auth.RefreshToken"

	hashRefreshToken := s.hasherManager.HashRefreshToken(refreshToken)

	tokenData, err := s.refreshRepo.GetByToken(ctx, s.db, hashRefreshToken)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrRefreshTokenNotFound) {
			return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidRefreshToken)
		}
		return domain.AuthTokens{}, fmt.Errorf("%s: %w: %v", op, errorsAPP.ErrInternalDB, err)
	}

	now := time.Now()

	if !tokenData.IsValid(now) {
		_ = s.refreshRepo.DeleteRefreshToken(ctx, s.db, hashRefreshToken)
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidRefreshToken)
	}

	user, err := s.repo.GetUserByID(ctx, s.db, tokenData.UserID)
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := user.CanLogin(); err != nil {
		return domain.AuthTokens{}, fmt.Errorf("%s: %w", op, err)
	}

	var newTokens domain.AuthTokens

	err = s.txmanager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {

		if err := s.refreshRepo.DeleteRefreshToken(ctx, q, hashRefreshToken); err != nil {
			return err
		}

		var err error
		newTokens, err = s.generateAndSaveTokens(ctx, q, user, now)
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
