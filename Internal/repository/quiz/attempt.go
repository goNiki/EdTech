package quiz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repositoryImpl) CreateAttempt(ctx context.Context, attempt *domain.QuizAttempt) (*domain.QuizAttempt, error) {
	const op = "repository.quiz.CreateAttempt"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	startedAt := attempt.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now()
		attempt.StartedAt = startedAt
	}

	query := `
		INSERT INTO quiz_attempts (
			quiz_id, 
			user_id, 
			score, 
			passed, 
			started_at, 
			completed_at, 
			created_at
		) VALUES (
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			$6, 
			NOW()
		)
		RETURNING id, created_at`

	var createdAt time.Time
	err := q.QueryRow(
		ctx,
		query,
		attempt.QuizID,
		attempt.UserID,
		attempt.Score,
		attempt.Passed,
		startedAt,
		attempt.CompletedAt,
	).Scan(
		&attempt.ID,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return attempt, nil
}

func (r *repositoryImpl) GetAttemptByID(ctx context.Context, id int64) (*domain.QuizAttempt, error) {
	const op = "repository.quiz.GetAttemptByID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
			id, 
			quiz_id, 
			user_id, 
			score, 
			passed, 
			COALESCE(draft_answers, '{}'::jsonb), 
			COALESCE(current_step, 1), 
			started_at, 
			completed_at 
		FROM quiz_attempts 
		WHERE id = $1`

	var attempt repomodels.QuizAttempt
	err := q.QueryRow(ctx, query, id).Scan(
		&attempt.ID,
		&attempt.QuizID,
		&attempt.UserID,
		&attempt.Score,
		&attempt.Passed,
		&attempt.DraftAnswers,
		&attempt.CurrentStep,
		&attempt.StartedAt,
		&attempt.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrAttemptNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.QuizAttemptToDomain(&attempt), nil
}

func (r *repositoryImpl) GetAttemptForUpdate(ctx context.Context, attemptID int64) (*domain.QuizAttempt, error) {
	const op = "repository.quiz.GetAttemptForUpdate"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
			id, 
			quiz_id, 
			user_id, 
			score, 
			passed, 
			started_at, 
			completed_at 
		FROM quiz_attempts 
		WHERE id = $1 
		FOR UPDATE`

	var attempt repomodels.QuizAttempt
	err := q.QueryRow(ctx, query, attemptID).Scan(
		&attempt.ID,
		&attempt.QuizID,
		&attempt.UserID,
		&attempt.Score,
		&attempt.Passed,
		&attempt.StartedAt,
		&attempt.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrAttemptNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.QuizAttemptToDomain(&attempt), nil
}

func (r *repositoryImpl) UpdateAttempt(ctx context.Context, attempt *domain.QuizAttempt) error {
	const op = "repository.quiz.UpdateAttempt"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		UPDATE quiz_attempts 
		SET score = $2, passed = $3, completed_at = $4 
		WHERE id = $1`

	tag, err := q.Exec(ctx, query, attempt.ID, attempt.Score, attempt.Passed, attempt.CompletedAt)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrAttemptNotFound)
	}

	return nil
}

func (r *repositoryImpl) CountUserAttempts(ctx context.Context, userID, quizID int64) (int, error) {
	const op = "repository.quiz.CountUserAttempts"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `SELECT COUNT(*) FROM quiz_attempts WHERE user_id = $1 AND quiz_id = $2`

	var count int
	err := q.QueryRow(ctx, query, userID, quizID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return count, nil
}

func (r *repositoryImpl) CountUserAttemptsForUpdate(ctx context.Context, userID, quizID int64) (int, error) {
	const op = "repository.quiz.CountUserAttemptsForUpdate"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `SELECT COUNT(*) FROM (SELECT id FROM quiz_attempts WHERE user_id = $1 AND quiz_id = $2 FOR UPDATE) AS locked_attempts`

	var count int
	err := q.QueryRow(ctx, query, userID, quizID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return count, nil
}

func (r *repositoryImpl) ListAttemptsForGrading(ctx context.Context, courseID int64, quizID *int64, limit, offset int) ([]domain.QuizAttempt, int64, error) {
	const op = "repository.quiz.ListAttemptsForGrading"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	countQuery := `
		SELECT COUNT(DISTINCT qa.id) 
		FROM quiz_attempts qa 
		JOIN quizzes q ON qa.quiz_id = q.id 
		WHERE q.course_id = $1 AND EXISTS (
			SELECT 1 FROM quiz_attempt_answers qaa WHERE qaa.attempt_id = qa.id AND qaa.is_correct IS NULL
		)`

	dataQuery := `
		SELECT 
			qa.id, 
			qa.quiz_id, 
			qa.user_id, 
			qa.score, 
			qa.passed, 
			qa.started_at, 
			qa.completed_at 
		FROM quiz_attempts qa 
		JOIN quizzes q ON qa.quiz_id = q.id 
		WHERE q.course_id = $1 AND EXISTS (
			SELECT 1 FROM quiz_attempt_answers qaa WHERE qaa.attempt_id = qa.id AND qaa.is_correct IS NULL
		)`

	countArgs := []interface{}{courseID}
	dataArgs := []interface{}{courseID}

	if quizID != nil && *quizID > 0 {
		countQuery += " AND qa.quiz_id = $2"
		dataQuery += " AND qa.quiz_id = $2"
		countArgs = append(countArgs, *quizID)
		dataArgs = append(dataArgs, *quizID)
	}

	var total int64
	err := q.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if limit > 0 {
		dataQuery += fmt.Sprintf(" ORDER BY qa.created_at DESC LIMIT $%d OFFSET $%d", len(dataArgs)+1, len(dataArgs)+2)
		dataArgs = append(dataArgs, limit, offset)
	} else {
		dataQuery += " ORDER BY qa.created_at DESC"
	}

	rows, err := q.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var attempts []domain.QuizAttempt
	for rows.Next() {
		var att repomodels.QuizAttempt
		if err := rows.Scan(
			&att.ID,
			&att.QuizID,
			&att.UserID,
			&att.Score,
			&att.Passed,
			&att.StartedAt,
			&att.CompletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}
		attempts = append(attempts, *repoconverter.QuizAttemptToDomain(&att))
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return attempts, total, nil
}


func (r *repositoryImpl) ListUserAttemptsByLessonID(ctx context.Context, userID, lessonID int64) ([]domain.QuizAttempt, error) {
	const op = "repository.quiz.ListUserAttemptsByLessonID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
			qa.id, 
			qa.quiz_id, 
			qa.user_id, 
			qa.score, 
			qa.passed, 
			qa.started_at, 
			qa.completed_at 
		FROM quiz_attempts qa
		JOIN quizzes q ON q.id = qa.quiz_id
		WHERE qa.user_id = $1 AND q.lesson_id = $2 AND q.deleted_at IS NULL
		ORDER BY qa.created_at ASC`

	rows, err := q.Query(ctx, query, userID, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var attempts []domain.QuizAttempt
	for rows.Next() {
		var att repomodels.QuizAttempt
		if err := rows.Scan(
			&att.ID,
			&att.QuizID,
			&att.UserID,
			&att.Score,
			&att.Passed,
			&att.StartedAt,
			&att.CompletedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}
		attempts = append(attempts, *repoconverter.QuizAttemptToDomain(&att))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return attempts, nil
}

func (r *repositoryImpl) GetBestScoreByLessonID(ctx context.Context, userID, lessonID int64) (int, error) {
	const op = "repository.quiz.GetBestScoreByLessonID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT COALESCE(MAX(qa.score), 0)
		FROM quiz_attempts qa
		JOIN quizzes q ON q.id = qa.quiz_id
		WHERE qa.user_id = $1 AND q.lesson_id = $2 AND q.deleted_at IS NULL AND qa.completed_at IS NOT NULL`

	var bestScore int
	err := q.QueryRow(ctx, query, userID, lessonID).Scan(&bestScore)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	return bestScore, nil
}


func (r *repositoryImpl) SaveAttemptDraft(ctx context.Context, attemptID, userID int64, currentStep int, draftAnswers []byte) error {
	const op = "repository.quiz.SaveAttemptDraft"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	if currentStep <= 0 {
		currentStep = 1
	}
	if len(draftAnswers) == 0 {
		draftAnswers = []byte("{}")
	}

	query := `
		UPDATE quiz_attempts 
		SET current_step = $3, draft_answers = $4 
		WHERE id = $1 AND user_id = $2 AND completed_at IS NULL`

	tag, err := q.Exec(ctx, query, attemptID, userID, currentStep, draftAnswers)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrAttemptNotFound)
	}
	return nil
}

func (r *repositoryImpl) GetActiveAttempt(ctx context.Context, userID, lessonID int64) (*domain.QuizAttempt, error) {
	const op = "repository.quiz.GetActiveAttempt"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
			qa.id, 
			qa.quiz_id, 
			qa.user_id, 
			qa.score, 
			qa.passed, 
			COALESCE(qa.draft_answers, '{}'::jsonb), 
			COALESCE(qa.current_step, 1), 
			qa.started_at, 
			qa.completed_at 
		FROM quiz_attempts qa
		JOIN quizzes q ON q.id = qa.quiz_id
		WHERE qa.user_id = $1 AND q.lesson_id = $2 AND q.deleted_at IS NULL AND qa.completed_at IS NULL
		ORDER BY qa.started_at DESC
		LIMIT 1`

	var attempt repomodels.QuizAttempt
	err := q.QueryRow(ctx, query, userID, lessonID).Scan(
		&attempt.ID,
		&attempt.QuizID,
		&attempt.UserID,
		&attempt.Score,
		&attempt.Passed,
		&attempt.DraftAnswers,
		&attempt.CurrentStep,
		&attempt.StartedAt,
		&attempt.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.QuizAttemptToDomain(&attempt), nil
}
