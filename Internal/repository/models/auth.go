package models

import "time"

type User struct {
	ID            int64
	Email         string
	PasswordHash  string
	Username      string
	FirstName     *string
	LastName      *string
	AvatarURL     *string
	Bio           *string
	Headline      *string
	Preferences   []byte
	Role          string
	EmailVerified bool
	IsActive      bool
	IsBanned      bool
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type CreateUser struct {
	Email        string
	PasswordHash string
	Username     string
	Role         string
}
