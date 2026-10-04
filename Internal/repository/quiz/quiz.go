package quiz

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repositoryImpl) CreateQuiz(ctx context.Context, quiz *domain.Quiz) (*domain.Quiz, error) {
	const op = "repository.quiz.CreateQuiz"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		INSERT INTO quizzes (
			lesson_id, 
			course_id, 
			title, 
			description, 
			passing_score, 
			max_attempts, 
			time_limit, 
			type, 
			created_at, 
			updated_at
		) VALUES (
			$1, 
			(SELECT course_id FROM lessons WHERE id = $1), 
			$2, 
			$3, 
			$4, 
			$5, 
			$6, 
			'single_choice', 
			NOW(), 
			NOW()
		)
		RETURNING id, created_at, updated_at`

	err := q.QueryRow(
		ctx,
		query,
		quiz.LessonID,
		quiz.Title,
		quiz.Description,
		quiz.PassingScor,
		quiz.MaxAttempts,
		quiz.TimeLimit,
	).Scan(
		&quiz.ID,
		&quiz.CreatedAt,
		&quiz.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return quiz, nil
}

func (r *repositoryImpl) GetQuizByID(ctx context.Context, id int64) (*domain.Quiz, error) {
	const op = "repository.quiz.GetQuizByID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
			id, 
			lesson_id, 
			title, 
			description, 
			passing_score, 
			max_attempts, 
			time_limit, 
			created_at, 
			updated_at, 
			deleted_at 
		FROM quizzes 
		WHERE id = $1 AND deleted_at IS NULL`

	var qz repomodels.Quiz
	err := q.QueryRow(ctx, query, id).Scan(
		&qz.ID,
		&qz.LessonID,
		&qz.Title,
		&qz.Description,
		&qz.PassingScor,
		&qz.MaxAttempts,
		&qz.TimeLimit,
		&qz.CreatedAt,
		&qz.UpdatedAt,
		&qz.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrQuizNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.QuizToDomain(&qz), nil
}

func (r *repositoryImpl) GetQuizTotalPoints(ctx context.Context, quizID int64) (int, error) {
	q := txmanager.GetQueryExecutor(ctx, r.Pool)
	var maxPoints int
	err := q.QueryRow(ctx, "SELECT COALESCE(SUM(points), 1) FROM quiz_questions WHERE quiz_id = $1", quizID).Scan(&maxPoints)
	if err != nil {
		return 1, err
	}
	if maxPoints == 0 {
		return 1, nil
	}
	return maxPoints, nil
}
