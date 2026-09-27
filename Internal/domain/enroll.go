package domain

type SelfEnrollRequest struct {
	UserID   int64
	CourseID int64
}

type TeacherEnrollRequest struct {
	TeacherID    int64
	CourseID     int64
	Role         string
	TargetEmail  string
	TargetUserID *int64
}
