package quiz

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repositoryImpl) CreateBatchAnswers(ctx context.Context, answers []domain.QuizAttemptAnswer) error {
	const op = "repository.quiz.CreateBatchAnswers"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	if len(answers) == 0 {
		return nil
	}

	query := `
		INSERT INTO quiz_attempt_answers (
			attempt_id, 
			question_id, 
			answer_id, 
			user_answer, 
			is_correct, 
			created_at
		) VALUES `

	args := make([]interface{}, 0, len(answers)*5)
	for i, ans := range answers {
		if i > 0 {
			query += ", "
		}
		offset := i * 5
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, NOW())", offset+1, offset+2, offset+3, offset+4, offset+5)
		args = append(args, ans.AttemptID, ans.QuestionID, ans.AnswerID, ans.TextValue, ans.IsCorrect)
	}

	_, err := q.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return nil
}

func (r *repositoryImpl) UpdateAttemptAnswer(ctx context.Context, answerID int64, points int, feedback *string, isCorrect bool) error {
	const op = "repository.quiz.UpdateAttemptAnswer"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		UPDATE quiz_attempt_answers 
		SET is_correct = $2, points = $3, feedback = $4, graded_at = NOW() 
		WHERE id = $1`

	tag, err := q.Exec(ctx, query, answerID, isCorrect, points, feedback)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrAnswerNotFound)
	}

	return nil
}

func (r *repositoryImpl) CountUngradedAnswers(ctx context.Context, attemptID int64) (int, error) {
	const op = "repository.quiz.CountUngradedAnswers"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `SELECT COUNT(*) FROM quiz_attempt_answers WHERE attempt_id = $1 AND is_correct IS NULL`

	var count int
	err := q.QueryRow(ctx, query, attemptID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return count, nil
}

func (r *repositoryImpl) GetAnswerPointsAndCorrectness(ctx context.Context, answerID int64) (bool, int, error) {
	const op = "repository.quiz.GetAnswerPointsAndCorrectness"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT qa.is_correct, qq.points 
		FROM quiz_answers qa 
		JOIN quiz_questions qq ON qa.question_id = qq.id 
		WHERE qa.id = $1`

	var isCorrect bool
	var points int
	err := q.QueryRow(ctx, query, answerID).Scan(&isCorrect, &points)
	if err != nil {
		return false, 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return isCorrect, points, nil
}

func (r *repositoryImpl) UpdateLessonProgressAfterQuiz(ctx context.Context, userID, lessonID int64, score int) error {
	const op = "repository.quiz.UpdateLessonProgressAfterQuiz"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		INSERT INTO lesson_progress (
			user_id, lesson_id, course_id, status, score, completed_at, last_accessed_at
		) VALUES (
			$1, $2, (SELECT course_id FROM lessons WHERE id = $2), 'completed', $3, NOW(), NOW()
		) ON CONFLICT (user_id, lesson_id) 
		DO UPDATE SET status = 'completed', score = $3, completed_at = NOW(), last_accessed_at = NOW()`

	_, err := q.Exec(ctx, query, userID, lessonID, score)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return nil
}

func (r *repositoryImpl) SumAttemptPoints(ctx context.Context, attemptID int64) (int, error) {
	const op = "repository.quiz.SumAttemptPoints"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `SELECT COALESCE(COUNT(*), 0) FROM quiz_attempt_answers WHERE attempt_id = $1 AND is_correct = TRUE`

	var totalPoints int
	err := q.QueryRow(ctx, query, attemptID).Scan(&totalPoints)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return totalPoints, nil
}
