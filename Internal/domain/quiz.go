package domain

import "time"

type QuizType string

const (
	QuizTypeSingleChoice QuizType = "single_choice"
	QuizTypeMultiChoice  QuizType = "multiple_choice"
	QuizTypeOpenText     QuizType = "open_text"
)

type Quiz struct {
	ID          int64      `db:"id"`
	LessonID    int64      `db:"lesson_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	PassingScor int        `db:"passing_score"`
	MaxAttempts *int       `db:"max_attempts"`
	TimeLimit   *int       `db:"time_limit"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type QuizQuestion struct {
	ID       int64    `db:"id"`
	QuizID   int64    `db:"quiz_id"`
	Type     QuizType `db:"type"`
	Text     string   `db:"question_text"`
	Points   int      `db:"points"`
	Position int      `db:"position"`
}

type QuizAnswer struct {
	ID         int64  `db:"id"`
	QuestionID int64  `db:"question_id"`
	Text       string `db:"answer_text"`
	IsCorrect  bool   `db:"is_correct"`
	Explain    string `db:"explanation"`
}

type QuizAttempt struct {
	ID          int64      `db:"id"`
	QuizID      int64      `db:"quiz_id"`
	UserID      int64      `db:"user_id"`
	Score       int        `db:"score"`
	Passed      bool       `db:"passed"`
	StartedAt   time.Time  `db:"started_at"`
	CompletedAt *time.Time `db:"completed_at"`
}

type QuizAttemptAnswer struct {
	ID         int64  `db:"id"`
	AttemptID  int64  `db:"attempt_id"`
	QuestionID int64  `db:"question_id"`
	AnswerID   *int64 `db:"answer_id"`
	TextValue  string `db:"text_value"`
	IsCorrect  *bool  `db:"is_correct"`
	Points     int    `db:"points_awarded"`
}
