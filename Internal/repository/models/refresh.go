package models

import "time"

type RefreshTokenData struct {
	ID           int64
	RefreshToken string
	UserID       int64
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
