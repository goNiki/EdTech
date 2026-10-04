package domain

import "time"

type ProgressStatus string

const (
	ProgressStatusNotStarted ProgressStatus = "not_started"
	ProgressStatusInProgress ProgressStatus = "in_progress"
	ProgressStatusCompleted  ProgressStatus = "completed"
)

type LessonProgress struct {
	ID          int64
	UserID      int64
	LessonID    int64
	CourseID    int64
	Status      ProgressStatus
	Score       *int
	TimeSpent   int
	LastPos     int
	StartedAt   *time.Time
	CompletedAt *time.Time
	UpdatedAt   time.Time
	Submissions []LessonSubmissionDetail
}

type CourseProgress struct {
	ID             int64
	UserID         int64
	CourseID       int64
	CompletedLess  int
	TotalLessons   int
	Percent        int
	TotalWatchTime int
	AverageScore   *float64
	StartedAt      *time.Time
	LastAccessedAt time.Time
	CompletedAt    *time.Time
}

type UpdateProgressInput struct {
	Status       ProgressStatus `json:"status"`
	TimeSpent    int            `json:"time_spent"`
	LastPosition int            `json:"last_position"`
}

type EssaySubmission struct {
	QuestionText string
	AnswerText   string
	MaxPoints    int
}

type LessonSubmissionDetail struct {
	QuestionText    string  `json:"question_text"`
	StudentAnswer   string  `json:"student_answer"`
	PointsAwarded   int     `json:"points_awarded"`
	MaxPoints       int     `json:"max_points"`
	TeacherFeedback *string `json:"teacher_feedback"`
	IsGraded        bool    `json:"is_graded"`
}

type CompleteLessonInput struct {
	Score       *int
	TimeSpent   int
	Answers     []LessonAnswerSubmission
	Essays      []EssaySubmission
	AttemptID   *int64
	IsAbandoned bool
}

type LessonCompletionResult struct {
	LessonID       int64                            `json:"lesson_id"`
	Status         string                           `json:"status"`
	Score          int                              `json:"score"`
	EarnedPoints   int                              `json:"earned_points"`
	TotalMaxPoints int                              `json:"total_max_points"`
	IsPassed       bool                             `json:"is_passed"`
	Results        map[string]BlockValidationResult `json:"results,omitempty"`
}

type BlockValidationResult struct {
	IsCorrect     bool   `json:"is_correct"`
	Feedback      string `json:"feedback,omitempty"`
	CorrectAnswer any    `json:"correct_answer,omitempty"`
}

type LessonAnswerSubmission struct {
	BlockID string `json:"block_id"`
	Answer  any    `json:"answer"`
}

type LessonAttemptsSummary struct {
	LessonID            int64
	TotalAttemptsMade   int
	MaxAttemptsAllowed  int
	CanStartNewAttempt  bool
	BestScore           int
	BestScorePercentage int
	IsPassed            bool
	PassingThreshold    int
	LastAttempt         *LessonAttemptItem
	AttemptsHistory     []LessonAttemptItem
}

type LessonAttemptItem struct {
	AttemptID   int64
	Score       int
	SubmittedAt time.Time
}

type StartAttemptResult struct {
	AttemptID int64
	StartedAt time.Time
}


