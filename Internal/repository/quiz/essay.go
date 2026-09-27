package quiz

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
)

func (r *repositoryImpl) SaveEssaySubmission(ctx context.Context, q db.QueryExecutor, userID, courseID, lessonID int64, essay domain.EssaySubmission) error {
	const op = "repository.quiz.SaveEssaySubmission"

	if essay.AnswerText == "" {
		return nil
	}

	// 1. Ensure Quiz exists for the lesson
	var quizID int64
	qQuery := `SELECT id FROM quizzes WHERE lesson_id = $1 AND type = 'essay' AND deleted_at IS NULL LIMIT 1`
	if err := q.QueryRow(ctx, qQuery, lessonID).Scan(&quizID); err != nil {
		maxPts := essay.MaxPoints
		if maxPts <= 0 {
			maxPts = 25
		}
		insertQuiz := `
			INSERT INTO quizzes (lesson_id, course_id, title, type, points) 
			VALUES ($1, $2, $3, 'essay', $4) 
			RETURNING id`
		if iErr := q.QueryRow(ctx, insertQuiz, lessonID, courseID, "Задание с развернутым ответом", maxPts).Scan(&quizID); iErr != nil {
			return fmt.Errorf("%s: create essay quiz: %w: %w", op, errorsAPP.ErrInternalDB, iErr)
		}
	}

	// 2. Ensure Question exists for the quiz
	var questionID int64
	qqQuery := `SELECT id FROM quiz_questions WHERE quiz_id = $1 LIMIT 1`
	if err := q.QueryRow(ctx, qqQuery, quizID).Scan(&questionID); err != nil {
		qText := essay.QuestionText
		if qText == "" {
			qText = "Развернутый ответ на вопрос"
		}
		insertQQ := `
			INSERT INTO quiz_questions (quiz_id, question_text, position) 
			VALUES ($1, $2, 1) 
			RETURNING id`
		if iErr := q.QueryRow(ctx, insertQQ, quizID, qText).Scan(&questionID); iErr != nil {
			return fmt.Errorf("%s: create essay question: %w: %w", op, errorsAPP.ErrInternalDB, iErr)
		}
	}

	// 3. Create Quiz Attempt
	var attemptID int64
	insAttempt := `
		INSERT INTO quiz_attempts (quiz_id, user_id, score, passed, started_at, completed_at) 
		VALUES ($1, $2, 0, false, NOW(), NOW()) 
		RETURNING id`
	if aErr := q.QueryRow(ctx, insAttempt, quizID, userID).Scan(&attemptID); aErr != nil {
		return fmt.Errorf("%s: create essay attempt: %w: %w", op, errorsAPP.ErrInternalDB, aErr)
	}

	// 4. Create Quiz Attempt Answer
	insAnswer := `
		INSERT INTO quiz_attempt_answers (attempt_id, question_id, user_answer, is_correct, points, created_at)
		VALUES ($1, $2, $3, NULL, 0, NOW())`
	if _, ansErr := q.Exec(ctx, insAnswer, attemptID, questionID, essay.AnswerText); ansErr != nil {
		return fmt.Errorf("%s: save essay answer: %w: %w", op, errorsAPP.ErrInternalDB, ansErr)
	}

	return nil
}
