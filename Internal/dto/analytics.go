package dto

import "time"

type CourseAnalyticsResponse struct {
	TotalStudents         int     `json:"total_students"`
	AvgProgressPercent    float64 `json:"avg_progress_percent"`
	AvgScore              float64 `json:"avg_score"`
	PendingHomeworksCount int     `json:"pending_homeworks_count"`
}

type PendingHomeworkDTO struct {
	AttemptID       int64     `json:"attempt_id"`
	AnswerID        int64     `json:"answer_id"`
	StudentID       int64     `json:"student_id"`
	StudentName     string    `json:"student_name"`
	StudentEmail    string    `json:"student_email"`
	StudentUsername string    `json:"student_username"`
	CourseID        int64     `json:"course_id"`
	CourseTitle     string    `json:"course_title"`
	LessonID        int64     `json:"lesson_id"`
	LessonTitle     string    `json:"lesson_title"`
	QuestionID      int64     `json:"question_id"`
	QuestionText    string    `json:"question_text"`
	StudentAnswer   string    `json:"student_answer"`
	AttachmentURL   *string   `json:"attachment_url,omitempty"`
	MaxPoints       int       `json:"max_points"`
	Rubric          *string   `json:"rubric,omitempty"`
	SubmittedAt     time.Time `json:"submitted_at"`
}

type PaginatedPendingHomeworksResponse struct {
	Items    []PendingHomeworkDTO `json:"items"`
	Total    int64                `json:"total"`
	Page     int64                `json:"page"`
	PageSize int64                `json:"page_size"`
}

type CoursePendingSummaryDTO struct {
	CourseID     int64  `json:"course_id"`
	CourseTitle  string `json:"course_title"`
	PendingCount int64  `json:"pending_count"`
}

type PaginatedTeacherPendingHomeworksResponse struct {
	Items          []PendingHomeworkDTO      `json:"items"`
	Total          int64                     `json:"total"`
	Page           int64                     `json:"page"`
	PageSize       int64                     `json:"page_size"`
	CoursesSummary []CoursePendingSummaryDTO `json:"courses_summary"`
}

type StudentLessonLogDTO struct {
	LessonID       int64      `json:"lesson_id"`
	LessonTitle    string     `json:"lesson_title"`
	LessonType     string     `json:"lesson_type"`
	Status         string     `json:"status"`
	Score          *int       `json:"score,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
}

type StudentAttemptAnswerDetailDTO struct {
	QuestionID    int64   `json:"question_id"`
	QuestionText  string  `json:"question_text"`
	ChosenAnswer  string  `json:"chosen_answer"`
	IsCorrect     *bool   `json:"is_correct,omitempty"`
	Points        int     `json:"points"`
	MaxPoints     int     `json:"max_points"`
	Feedback      *string `json:"feedback,omitempty"`
	AttachmentURL *string `json:"attachment_url,omitempty"`
	Explanation   string  `json:"explanation,omitempty"`
}

type StudentTestAttemptDTO struct {
	AttemptID     int64                           `json:"attempt_id"`
	QuizID        int64                           `json:"quiz_id"`
	QuizTitle     string                          `json:"quiz_title"`
	AttemptNumber int                             `json:"attempt_number"`
	Score         int                             `json:"score"`
	Passed        bool                            `json:"passed"`
	StartedAt     time.Time                       `json:"started_at"`
	CompletedAt   *time.Time                      `json:"completed_at,omitempty"`
	Answers       []StudentAttemptAnswerDetailDTO `json:"answers"`
}

type StudentDrilldownResponse struct {
	StudentID       int64                   `json:"student_id"`
	StudentName     string                  `json:"student_name"`
	StudentEmail    string                  `json:"student_email"`
	StudentUsername string                  `json:"student_username"`
	AvatarURL       *string                 `json:"avatar_url,omitempty"`
	Bio             *string                 `json:"bio,omitempty"`
	CourseID        int64                   `json:"course_id"`
	OverallProgress float64                 `json:"overall_progress"`
	AvgScore        float64                 `json:"avg_score"`
	LessonLogs      []StudentLessonLogDTO   `json:"lesson_logs"`
	TestAttempts    []StudentTestAttemptDTO `json:"test_attempts"`
}
