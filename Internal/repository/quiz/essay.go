package quiz

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repositoryImpl) SaveEssaySubmission(ctx context.Context, userID, courseID, lessonID int64, essay domain.EssaySubmission) error {
	const op = "repository.quiz.SaveEssaySubmission"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

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

func (r *repositoryImpl) GetLessonSubmissions(ctx context.Context, userID, lessonID int64) ([]domain.LessonSubmissionDetail, error) {
	const op = "repository.quiz.GetLessonSubmissions"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
			COALESCE(qq.question_text, ''),
			COALESCE(qaa.user_answer, ''),
			COALESCE(qaa.points, 0),
			COALESCE(qz.points, 25),
			qaa.feedback,
			(qaa.is_correct IS NOT NULL) AS is_graded
		FROM quiz_attempt_answers qaa
		JOIN quiz_attempts qa ON qa.id = qaa.attempt_id
		JOIN quizzes qz ON qz.id = qa.quiz_id
		JOIN quiz_questions qq ON qq.id = qaa.question_id
		WHERE qz.lesson_id = $1 
		  AND qa.user_id = $2
		  AND qz.deleted_at IS NULL
		ORDER BY qaa.created_at ASC`

	rows, err := q.Query(ctx, query, lessonID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	submissions := make([]domain.LessonSubmissionDetail, 0)
	for rows.Next() {
		var sub domain.LessonSubmissionDetail
		if err := rows.Scan(
			&sub.QuestionText,
			&sub.StudentAnswer,
			&sub.PointsAwarded,
			&sub.MaxPoints,
			&sub.TeacherFeedback,
			&sub.IsGraded,
		); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}
		submissions = append(submissions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return submissions, nil
}

