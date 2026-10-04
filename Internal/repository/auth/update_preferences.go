package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repository) UpdatePreferences(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error) {
	const op = "repository.auth.UpdatePreferences"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	patch := make(map[string]string)
	if input.FontScale != nil {
		patch["font_scale"] = *input.FontScale
	}
	if input.ContentWidth != nil {
		patch["content_width"] = *input.ContentWidth
	}
	if input.LineHeight != nil {
		patch["line_height"] = *input.LineHeight
	}
	if input.ReadingTheme != nil {
		patch["reading_theme"] = *input.ReadingTheme
	}

	patchJSON, err := json.Marshal(patch)
	if err != nil {
		return domain.UserPreferences{}, fmt.Errorf("%s: marshal patch: %w", op, err)
	}

	query := `
		UPDATE users
		SET preferences = COALESCE(preferences, '{"font_scale": "medium", "content_width": "wide", "line_height": "normal", "reading_theme": "system"}'::jsonb) || $2::jsonb,
		    updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING preferences
	`

	var rawJSON []byte
	if err := q.QueryRow(ctx, query, userID, patchJSON).Scan(&rawJSON); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserPreferences{}, errorsAPP.ErrUserNotFound
		}
		return domain.UserPreferences{}, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	result := domain.DefaultUserPreferences()
	if len(rawJSON) > 0 {
		if err := json.Unmarshal(rawJSON, &result); err != nil {
			return domain.UserPreferences{}, fmt.Errorf("%s: unmarshal preferences: %w", op, err)
		}
	}

	return result, nil
}
