package domain

import "time"

type Role string

const (
	RoleStudent Role = "user"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Username     string
	Role         Role
	CreateAt     time.Time
	UpdateAt     time.Time
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
