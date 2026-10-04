package domain

import "time"

type CourseAnalyticsSummary struct {
	TotalStudents         int
	AvgProgressPercent    float64
	AvgScore              float64
	PendingHomeworksCount int
}

type CourseStudentItem struct {
	UserID              int64
	FirstName           string
	LastName            string
	Username            string
	Email               string
	AvatarURL           *string
	Role                string
	EnrolledAt          time.Time
	ProgressPercentage  float64
	CompletedLessons    int
	TotalLessons        int
	AverageScore        float64
	HasPendingHomeworks bool
}

type PendingHomeworkItem struct {
	AttemptID       int64
	AnswerID        int64
	StudentID       int64
	StudentName     string
	StudentEmail    string
	StudentUsername string
	CourseID        int64
	CourseTitle     string
	LessonID        int64
	LessonTitle     string
	QuestionID      int64
	QuestionText    string
	StudentAnswer   string
	AttachmentURL   *string
	MaxPoints       int
	Rubric          *string
	SubmittedAt     time.Time
}

type CoursePendingSummaryItem struct {
	CourseID     int64
	CourseTitle  string
	PendingCount int64
}

type TeacherPendingHomeworksResult struct {
	Items          []PendingHomeworkItem
	Total          int64
	CoursesSummary []CoursePendingSummaryItem
}

type StudentLessonLog struct {
	LessonID       int64
	LessonTitle    string
	LessonType     string
	Status         string
	Score          *int
	CompletedAt    *time.Time
	LastAccessedAt *time.Time
}

type StudentAttemptAnswerDetail struct {
	QuestionID    int64
	QuestionText  string
	ChosenAnswer  string
	CorrectAnswer string
	Points        int
	MaxPoints     int
	Feedback      *string
	AttachmentURL *string
	Explanation   string
	IsGraded      bool
	IsCorrect     *bool
}

type StudentTestAttemptLog struct {
	AttemptID     int64
	QuizID        int64
	QuizTitle     string
	AttemptNumber int
	Score         int
	Passed        bool
	StartedAt     time.Time
	CompletedAt   *time.Time
	Answers       []StudentAttemptAnswerDetail
}

type StudentDrilldownReport struct {
	Student         User
	CourseID        int64
	OverallProgress float64
	AvgScore        float64
	LessonLogs      []StudentLessonLog
	TestAttempts    []StudentTestAttemptLog
}
