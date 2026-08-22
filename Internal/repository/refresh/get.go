package refresh

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetByToken(ctx context.Context, q db.QueryExecutor, refreshToken string) (domain.RefreshTokenData, error) {
	const op = "repository.refresh.GetByToken"

	query := `
		SELECT 
			id,
			token,
			user_id,
			expires_at,
			created_at,
		FROM refresh_tokens
		WHERE token = $1
	`

	var tokenData repomodels.RefreshTokenData

	err := q.QueryRow(
		ctx,
		query,
		refreshToken,
	).Scan(
		&tokenData.ID,
		&tokenData.RefreshToken,
		&tokenData.UserID,
		&tokenData.ExpiresAt,
		&tokenData.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshTokenData{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrRefreshTokenNotFound)
		}

		return domain.RefreshTokenData{}, fmt.Errorf("%s: %w", op, err)
	}

	return repoconverter.RefresTokenToDomain(tokenData), nil

}
