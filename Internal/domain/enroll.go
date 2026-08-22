package domain

type EnrollUserRequest struct {
	UserEmail string
	CourseID  int64
	Role      string
}