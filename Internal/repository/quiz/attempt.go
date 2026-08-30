package quiz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repositoryImpl) CreateAttempt(ctx context.Context, q db.QueryExecutor, attempt *domain.QuizAttempt) (*domain.QuizAttempt, error) {
	const op = "repository.quiz.CreateAttempt"

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

func (r *repositoryImpl) GetAttemptByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.QuizAttempt, error) {
	const op = "repository.quiz.GetAttemptByID"

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
		WHERE id = $1`

	var attempt repomodels.QuizAttempt
	err := q.QueryRow(ctx, query, id).Scan(
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

func (r *repositoryImpl) GetAttemptForUpdate(ctx context.Context, q db.QueryExecutor, attemptID int64) (*domain.QuizAttempt, error) {
	const op = "repository.quiz.GetAttemptForUpdate"

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

func (r *repositoryImpl) UpdateAttempt(ctx context.Context, q db.QueryExecutor, attempt *domain.QuizAttempt) error {
	const op = "repository.quiz.UpdateAttempt"

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

func (r *repositoryImpl) CountUserAttempts(ctx context.Context, q db.QueryExecutor, userID, quizID int64) (int, error) {
	const op = "repository.quiz.CountUserAttempts"

	query := `SELECT COUNT(*) FROM quiz_attempts WHERE user_id = $1 AND quiz_id = $2`

	var count int
	err := q.QueryRow(ctx, query, userID, quizID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return count, nil
}

func (r *repositoryImpl) CountUserAttemptsForUpdate(ctx context.Context, q db.QueryExecutor, userID, quizID int64) (int, error) {
	const op = "repository.quiz.CountUserAttemptsForUpdate"

	query := `SELECT COUNT(*) FROM (SELECT id FROM quiz_attempts WHERE user_id = $1 AND quiz_id = $2 FOR UPDATE) AS locked_attempts`

	var count int
	err := q.QueryRow(ctx, query, userID, quizID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return count, nil
}

func (r *repositoryImpl) ListAttemptsForGrading(ctx context.Context, q db.QueryExecutor, courseID int64, quizID *int64, limit, offset int) ([]domain.QuizAttempt, int64, error) {
	const op = "repository.quiz.ListAttemptsForGrading"

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
