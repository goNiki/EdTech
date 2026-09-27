package progress

import (
	"context"
	"errors"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) CompleteLesson(ctx context.Context, userID int64, lessonID int64, score *int, essays []domain.EssaySubmission) error {
	const op = "service.progress.CompleteLesson"

	lesson, err := s.lessonRepo.GetLessonByID(ctx, s.db, lessonID)
	if err != nil {
		return fmt.Errorf("%s: get lesson: %w", op, err)
	}

	finalScore := 100
	if score != nil {
		finalScore = *score
	}

	err = s.progressRepo.UpdateLessonProgressStatus(ctx, s.db, userID, lessonID, domain.ProgressStatusCompleted)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrLessonProgressNotFound) {
			now := time.Now()
			lp := &domain.LessonProgress{
				UserID:      userID,
				LessonID:    lessonID,
				CourseID:    lesson.CourseID,
				Status:      domain.ProgressStatusCompleted,
				Score:       &finalScore,
				CompletedAt: &now,
				StartedAt:   &now,
				UpdatedAt:   now,
			}
			if cErr := s.progressRepo.CreateLessonProgress(ctx, s.db, lp); cErr != nil {
				return fmt.Errorf("%s: create lesson progress: %w", op, cErr)
			}
		} else {
			return fmt.Errorf("%s: update lesson progress: %w", op, err)
		}
	} else {
		updateScoreQuery := `UPDATE lesson_progress SET score = $1 WHERE user_id = $2 AND lesson_id = $3`
		_, _ = s.db.Exec(ctx, updateScoreQuery, finalScore, userID, lessonID)
	}

	// Handle Homework / Essay submissions
	if len(essays) > 0 {
		for _, es := range essays {
			if es.AnswerText == "" {
				continue
			}
			_ = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, tx db.QueryExecutor) error {
				// 1. Ensure Quiz exists
				var quizID int64
				qQuery := `SELECT id FROM quizzes WHERE lesson_id = $1 AND type = 'essay' AND deleted_at IS NULL LIMIT 1`
				qErr := tx.QueryRow(ctx, qQuery, lessonID).Scan(&quizID)
				if qErr != nil {
					maxPts := es.MaxPoints
					if maxPts <= 0 {
						maxPts = 25
					}
					insertQuiz := `INSERT INTO quizzes (lesson_id, course_id, title, type, points) VALUES ($1, $2, $3, 'essay', $4) RETURNING id`
					if iErr := tx.QueryRow(ctx, insertQuiz, lessonID, lesson.CourseID, "Задание с развернутым ответом", maxPts).Scan(&quizID); iErr != nil {
						return iErr
					}
				}

				// 2. Ensure Question exists
				var questionID int64
				qqQuery := `SELECT id FROM quiz_questions WHERE quiz_id = $1 LIMIT 1`
				qqErr := tx.QueryRow(ctx, qqQuery, quizID).Scan(&questionID)
				if qqErr != nil {
					qText := es.QuestionText
					if qText == "" {
						qText = "Развернутый ответ на вопрос"
					}
					insertQQ := `INSERT INTO quiz_questions (quiz_id, question_text, position) VALUES ($1, $2, 1) RETURNING id`
					if iErr := tx.QueryRow(ctx, insertQQ, quizID, qText).Scan(&questionID); iErr != nil {
						return iErr
					}
				}

				// 3. Create Quiz Attempt
				var attemptID int64
				insAttempt := `INSERT INTO quiz_attempts (quiz_id, user_id, score, passed, started_at, completed_at) 
				VALUES ($1, $2, 0, false, NOW(), NOW()) RETURNING id`
				if aErr := tx.QueryRow(ctx, insAttempt, quizID, userID).Scan(&attemptID); aErr != nil {
					return aErr
				}

				// 4. Create Quiz Attempt Answer
				insAnswer := `INSERT INTO quiz_attempt_answers (attempt_id, question_id, user_answer, is_correct, points, created_at)
				VALUES ($1, $2, $3, NULL, 0, NOW())`
				_, ansErr := tx.Exec(ctx, insAnswer, attemptID, questionID, es.AnswerText)
				return ansErr
			})
		}
	}

	// Recalculate Course Progress
	courseLessons, err := s.lessonRepo.GetLessonsByCourseID(ctx, s.db, lesson.CourseID)
	if err == nil && len(courseLessons) > 0 {
		allProgress, pErr := s.progressRepo.GetAllLessonProgressByCourse(ctx, s.db, userID, lesson.CourseID)
		if pErr == nil {
			completedCount := 0
			totalScore := 0
			scoreCount := 0
			for _, p := range allProgress {
				if p.Status == domain.ProgressStatusCompleted {
					completedCount++
					if p.Score != nil && *p.Score > 0 {
						totalScore += *p.Score
						scoreCount++
					}
				}
			}
			total := len(courseLessons)
			percent := (float64(completedCount) / float64(total)) * 100.0
			avgScore := 100.0
			if scoreCount > 0 {
				avgScore = float64(totalScore) / float64(scoreCount)
			}
			updateCP := `
				INSERT INTO course_progress (user_id, course_id, completed_lessons, total_lessons, progress_percentage, average_score, started_at, last_accessed_at, completed_at)
				VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), CASE WHEN $5 >= 100 THEN NOW() ELSE NULL END)
				ON CONFLICT (user_id, course_id) DO UPDATE
				SET completed_lessons = EXCLUDED.completed_lessons,
				    total_lessons = EXCLUDED.total_lessons,
				    progress_percentage = EXCLUDED.progress_percentage,
				    average_score = EXCLUDED.average_score,
				    last_accessed_at = NOW(),
				    completed_at = CASE WHEN EXCLUDED.progress_percentage >= 100 AND course_progress.completed_at IS NULL THEN NOW() ELSE course_progress.completed_at END
			`
			_, _ = s.db.Exec(ctx, updateCP, userID, lesson.CourseID, completedCount, total, percent, avgScore)
		}
	}

	return nil
}
