package dto

import "time"

type Quiz struct {
	ID          int64      `json:"id"`
	LessonID    int64      `json:"lesson_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	PassingScor int        `json:"passing_score"`
	MaxAttempts *int       `json:"max_attempts,omitempty"`
	TimeLimit   *int       `json:"time_limit,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type QuizQuestion struct {
	ID       int64        `json:"id"`
	QuizID   int64        `json:"quiz_id"`
	Type     string       `json:"type"` // single_choice, multiple_choice, open_text
	Text     string       `json:"question_text"`
	Points   int          `json:"points"`
	Position int          `json:"position"`
	Answers  []QuizAnswer `json:"answers,omitempty"`
}

type QuizAnswer struct {
	ID         int64   `json:"id"`
	QuestionID int64   `json:"question_id"`
	Text       string  `json:"answer_text"`
	IsCorrect  *bool   `json:"is_correct,omitempty"`
	Explain    *string `json:"explanation,omitempty"`
}

type CreateQuizRequest struct {
	LessonID    int64   `json:"lesson_id" validate:"required,gt=0"`
	Title       string  `json:"title" validate:"required,min=2,max=255"`
	Description string  `json:"description,omitempty" validate:"omitempty,max=2000"`
	PassingScor int     `json:"passing_score" validate:"required,min=0,max=100"`
	MaxAttempts *int    `json:"max_attempts,omitempty" validate:"omitempty,gt=0"`
	TimeLimit   *int    `json:"time_limit,omitempty" validate:"omitempty,gt=0"`
}

type AddQuestionRequest struct {
	Type     string            `json:"type" validate:"required,oneof=single_choice multiple_choice open_text"`
	Text     string            `json:"question_text" validate:"required,min=2,max=2000"`
	Points   int               `json:"points" validate:"required,gt=0"`
	Answers  []AddAnswerOption `json:"answers,omitempty" validate:"omitempty,dive"`
}

type AddAnswerOption struct {
	Text      string  `json:"answer_text" validate:"required,min=1,max=1000"`
	IsCorrect bool    `json:"is_correct"`
	Explain   *string `json:"explanation,omitempty" validate:"omitempty,max=1000"`
}

type SubmitAttemptRequest struct {
	Answers []SubmitAnswer `json:"answers" validate:"required,min=1,dive"`
}

type SubmitAnswer struct {
	QuestionID int64   `json:"question_id" validate:"required,gt=0"`
	AnswerID   *int64  `json:"answer_id,omitempty" validate:"omitempty,gt=0"`
	TextValue  *string `json:"text_value,omitempty" validate:"omitempty,max=5000"`
}
