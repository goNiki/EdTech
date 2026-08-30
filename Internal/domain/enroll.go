package domain

type SelfEnrollRequest struct {
	UserID   int64
	CourseID int64
}

type TeacherEnrollRequest struct {
	TeacherID   int64
	TargetEmail string
	CourseID    int64
	Role        string
}
