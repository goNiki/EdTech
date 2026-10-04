package domain

import (
	errorsAPP "edtech/pkg/errors"
	"time"
)

type QuizType string

const (
	QuizTypeSingleChoice QuizType = "single_choice"
	QuizTypeMultiChoice  QuizType = "multiple_choice"
	QuizTypeOpenText     QuizType = "open_text"
)

type Quiz struct {
	ID          int64
	LessonID    int64
	Title       string
	Description string
	PassingScor int
	MaxAttempts *int
	TimeLimit   *int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (q *Quiz) Validate() error {
	if q.Title == "" {
		return errorsAPP.ErrEmptyTitle
	}
	if q.LessonID <= 0 {
		return errorsAPP.ErrQuizValidation
	}
	if q.PassingScor < 0 || q.PassingScor > 100 {
		return errorsAPP.ErrQuizValidation
	}
	if q.MaxAttempts != nil && *q.MaxAttempts <= 0 {
		return errorsAPP.ErrQuizValidation
	}
	if q.TimeLimit != nil && *q.TimeLimit <= 0 {
		return errorsAPP.ErrQuizValidation
	}
	return nil
}

type QuizQuestion struct {
	ID       int64
	QuizID   int64
	Type     QuizType
	Text     string
	Points   int
	Position int
}

type QuizAnswer struct {
	ID         int64
	QuestionID int64
	Text       string
	IsCorrect  bool
	Explain    string
}

type QuizAttempt struct {
	ID           int64
	QuizID       int64
	UserID       int64
	Score        int
	Passed       bool
	NeedsGrading bool
	DraftAnswers map[string]any
	CurrentStep  int
	StartedAt    time.Time
	CompletedAt  *time.Time
}

type ActiveAttemptInfo struct {
	ID               int64          `json:"id"`
	StartedAt        time.Time      `json:"started_at"`
	TimeLimitMinutes int            `json:"time_limit_minutes"`
	RemainingSeconds int            `json:"remaining_seconds"`
	CurrentStep      int            `json:"current_step"`
	DraftAnswers     map[string]any `json:"draft_answers"`
}

type ActiveAttemptResult struct {
	HasActiveAttempt bool               `json:"has_active_attempt"`
	Attempt          *ActiveAttemptInfo `json:"attempt,omitempty"`
}

func (a *QuizAttempt) CalculateScore(correctPoints, totalQuestions, passingScore int) {
	if totalQuestions == 0 {
		totalQuestions = 1
	}
	a.Score = (correctPoints * 100) / totalQuestions
	if a.Score > 100 {
		a.Score = 100
	}
	a.Passed = a.Score >= passingScore
	a.NeedsGrading = false
	now := time.Now()
	a.CompletedAt = &now
}

type QuizAttemptAnswer struct {
	ID            int64
	AttemptID     int64
	QuestionID    int64
	AnswerID      *int64
	TextValue     string
	IsCorrect     *bool
	Points        int
	Feedback      *string
	AttachmentURL *string
	GradedBy      *int64
	GradedAt      *time.Time
	CreatedAt     time.Time
}
