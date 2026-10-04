package auth

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *service) UpdatePreferences(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error) {
	const op = "service.auth.UpdatePreferences"

	if err := input.Validate(); err != nil {
		return domain.UserPreferences{}, fmt.Errorf("%s: %w", op, err)
	}

	prefs, err := s.repo.UpdatePreferences(ctx, userID, input)
	if err != nil {
		return domain.UserPreferences{}, fmt.Errorf("%s: %w", op, err)
	}

	return prefs, nil
}
