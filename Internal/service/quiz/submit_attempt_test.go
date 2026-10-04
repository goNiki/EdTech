package quiz_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	"edtech/internal/service"
	"edtech/internal/service/quiz"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFullTx struct {
	committed  bool
	rolledBack bool
}

func (m *mockFullTx) Begin(ctx context.Context) (pgx.Tx, error) { return m, nil }
func (m *mockFullTx) Commit(ctx context.Context) error {
	m.committed = true
	return nil
}
func (m *mockFullTx) Rollback(ctx context.Context) error {
	m.rolledBack = true
	return nil
}
func (m *mockFullTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *mockFullTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *mockFullTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row { return nil }
func (m *mockFullTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (m *mockFullTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults { return nil }
func (m *mockFullTx) LargeObjects() pgx.LargeObjects                             { return pgx.LargeObjects{} }
func (m *mockFullTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (m *mockFullTx) Conn() *pgx.Conn { return nil }

type mockPoolSubmit struct {
	tx *mockFullTx
}

func (p *mockPoolSubmit) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.tx, nil
}

type mockQuizRepoForSubmit struct {
	repository.QuizRepository
	attempt          *domain.QuizAttempt
	quiz             *domain.Quiz
	batchAnswers     []domain.QuizAttemptAnswer
	updateAttemptErr error
}

func (m *mockQuizRepoForSubmit) GetAttemptForUpdate(ctx context.Context, attemptID int64) (*domain.QuizAttempt, error) {
	if m.attempt == nil {
		return nil, errorsAPP.ErrAttemptNotFound
	}
	return m.attempt, nil
}

func (m *mockQuizRepoForSubmit) GetQuizByID(ctx context.Context, id int64) (*domain.Quiz, error) {
	if m.quiz == nil {
		return nil, errorsAPP.ErrQuizNotFound
	}
	return m.quiz, nil
}

func (m *mockQuizRepoForSubmit) GetAnswerPointsAndCorrectness(ctx context.Context, answerID int64) (bool, int, error) {
	return true, 10, nil
}

func (m *mockQuizRepoForSubmit) CreateBatchAnswers(ctx context.Context, answers []domain.QuizAttemptAnswer) error {
	m.batchAnswers = answers
	return nil
}

func (m *mockQuizRepoForSubmit) GetQuizTotalPoints(ctx context.Context, quizID int64) (int, error) {
	return 10, nil
}

func (m *mockQuizRepoForSubmit) UpdateAttempt(ctx context.Context, attempt *domain.QuizAttempt) error {
	if m.updateAttemptErr != nil {
		return m.updateAttemptErr
	}
	return nil
}

func (m *mockQuizRepoForSubmit) UpdateLessonProgressAfterQuiz(ctx context.Context, userID, lessonID int64, score int) error {
	return nil
}

type mockProgressServiceForSubmit struct {
	service.ProgressServices
	completeErr error
	calledWith  *domain.CompleteLessonInput
}

func (m *mockProgressServiceForSubmit) CompleteLesson(ctx context.Context, userID int64, lessonID int64, input domain.CompleteLessonInput) (*domain.LessonCompletionResult, error) {
	m.calledWith = &input
	if m.completeErr != nil {
		return nil, m.completeErr
	}
	return &domain.LessonCompletionResult{
		LessonID: lessonID,
		Score:    *input.Score,
		IsPassed: true,
	}, nil
}

func TestSubmitAttempt_SuccessAtomicCommit(t *testing.T) {
	ctx := context.Background()
	tx := &mockFullTx{}
	pool := &mockPoolSubmit{tx: tx}
	tm := txmanager.NewTxManagerWithPool(pool)

	qRepo := &mockQuizRepoForSubmit{
		attempt: &domain.QuizAttempt{
			ID:        1,
			QuizID:    10,
			UserID:    42,
			StartedAt: time.Now().Add(-5 * time.Minute),
		},
		quiz: &domain.Quiz{
			ID:          10,
			LessonID:    100,
			PassingScor: 70,
		},
	}

	progressSvc := &mockProgressServiceForSubmit{}

	svc := quiz.NewQuizService(qRepo, nil, nil, nil, progressSvc, tm, nil)

	ansID := int64(99)
	answers := []domain.QuizAttemptAnswer{
		{
			QuestionID: 5,
			AnswerID:   &ansID,
		},
	}

	res, err := svc.SubmitAttempt(ctx, 42, 1, answers)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.True(t, res.Passed)
	assert.Equal(t, 100, res.Score)
	assert.True(t, tx.committed, "transaction must be committed")
	assert.False(t, tx.rolledBack, "transaction must not be rolled back")
	assert.NotNil(t, progressSvc.calledWith)
}

func TestSubmitAttempt_CompleteLessonFails_AtomicRollback(t *testing.T) {
	ctx := context.Background()
	tx := &mockFullTx{}
	pool := &mockPoolSubmit{tx: tx}
	tm := txmanager.NewTxManagerWithPool(pool)

	qRepo := &mockQuizRepoForSubmit{
		attempt: &domain.QuizAttempt{
			ID:        1,
			QuizID:    10,
			UserID:    42,
			StartedAt: time.Now().Add(-5 * time.Minute),
		},
		quiz: &domain.Quiz{
			ID:          10,
			LessonID:    100,
			PassingScor: 70,
		},
	}

	expectedErr := errors.New("db connection failure in progress service")
	progressSvc := &mockProgressServiceForSubmit{
		completeErr: expectedErr,
	}

	svc := quiz.NewQuizService(qRepo, nil, nil, nil, progressSvc, tm, nil)

	ansID := int64(99)
	answers := []domain.QuizAttemptAnswer{
		{
			QuestionID: 5,
			AnswerID:   &ansID,
		},
	}

	_, err := svc.SubmitAttempt(ctx, 42, 1, answers)
	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())

	assert.False(t, tx.committed, "transaction must NOT be committed when CompleteLesson fails")
	assert.True(t, tx.rolledBack, "transaction MUST be rolled back for atomic consistency")
}
