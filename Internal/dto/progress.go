package dto

import "time"

type LessonProgress struct {
	ID          int64                       `json:"id,omitempty"`
	LessonID    int64                       `json:"lesson_id"`
	Status      string                      `json:"status"` // not_started, in_progress, completed
	Score       *int                        `json:"score"`
	TimeSpent   int                         `json:"time_spent"`
	LastPos     int                         `json:"last_position"`
	StartedAt   *time.Time                  `json:"started_at,omitempty"`
	CompletedAt *time.Time                  `json:"completed_at,omitempty"`
	Submissions []LessonSubmissionDetailDTO `json:"submissions"`
}

type LessonSubmissionDetailDTO struct {
	QuestionText    string  `json:"question_text"`
	StudentAnswer   string  `json:"student_answer"`
	PointsAwarded   int     `json:"points_awarded"`
	MaxPoints       int     `json:"max_points"`
	TeacherFeedback *string `json:"teacher_feedback"`
	IsGraded        bool    `json:"is_graded"`
}

type CourseProgress struct {
	ID             int64      `json:"id"`
	CourseID       int64      `json:"course_id"`
	Status         string     `json:"status"`
	Percent        int        `json:"completion_percentage"`
	CompletedLess  int        `json:"completed_lessons"`
	TotalLessons   int        `json:"total_lessons"`
	AverageScore   *float64   `json:"average_score,omitempty"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	LastAccessedAt time.Time  `json:"last_accessed_at"`
}

type UpdateLessonProgressRequest struct {
	Status    string `json:"status" validate:"omitempty,oneof=not_started in_progress completed"`
	TimeSpent int    `json:"time_spent" validate:"gte=0"`
	LastPos   int    `json:"last_position" validate:"gte=0"`
}

type CompleteLessonRequest struct {
	Score     *int       `json:"score,omitempty"`
	TimeSpent int        `json:"time_spent,omitempty"`
	Essays    []EssayDTO `json:"essays,omitempty"`
}

type EssayDTO struct {
	QuestionText string `json:"question_text"`
	AnswerText   string `json:"answer_text"`
	MaxPoints    int    `json:"max_points"`
}
