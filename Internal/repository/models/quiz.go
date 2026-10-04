package models

import "time"

type Quiz struct {
	ID          int64
	LessonID    int64
	Title       string
	Description *string
	PassingScor int
	MaxAttempts *int
	TimeLimit   *int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type QuizQuestion struct {
	ID       int64
	QuizID   int64
	Type     string
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
	DraftAnswers []byte
	CurrentStep  int
	StartedAt    time.Time
	CompletedAt  *time.Time
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
