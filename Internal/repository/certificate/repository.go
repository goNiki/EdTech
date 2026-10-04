package certificate

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	Pool *pgxpool.Pool
}

func NewCertificateRepository(pool *pgxpool.Pool) repository.CertificateRepository {
	return &repo{Pool: pool}
}

func (r *repo) CreateCertificate(ctx context.Context, cert *domain.Certificate) (*domain.Certificate, error) {
	const op = "repository.certificate.CreateCertificate"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		INSERT INTO certificates (certificate_code, user_id, course_id, student_name, course_title, final_score, issued_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (user_id, course_id)
		DO UPDATE SET final_score = EXCLUDED.final_score
		RETURNING id, certificate_code, user_id, course_id, student_name, course_title, final_score, issued_at
	`

	var res domain.Certificate
	err := q.QueryRow(ctx, query, cert.CertificateCode, cert.UserID, cert.CourseID, cert.StudentName, cert.CourseTitle, cert.FinalScore).
		Scan(&res.ID, &res.CertificateCode, &res.UserID, &res.CourseID, &res.StudentName, &res.CourseTitle, &res.FinalScore, &res.IssuedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}

func (r *repo) GetCertificateByCode(ctx context.Context, code string) (*domain.Certificate, error) {
	const op = "repository.certificate.GetCertificateByCode"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT id, certificate_code, user_id, course_id, student_name, course_title, final_score, issued_at
		FROM certificates
		WHERE certificate_code = $1
	`

	var res domain.Certificate
	err := q.QueryRow(ctx, query, code).
		Scan(&res.ID, &res.CertificateCode, &res.UserID, &res.CourseID, &res.StudentName, &res.CourseTitle, &res.FinalScore, &res.IssuedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCertificateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}

func (r *repo) GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int64) (*domain.Certificate, error) {
	const op = "repository.certificate.GetCertificateByUserAndCourse"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT id, certificate_code, user_id, course_id, student_name, course_title, final_score, issued_at
		FROM certificates
		WHERE user_id = $1 AND course_id = $2
	`

	var res domain.Certificate
	err := q.QueryRow(ctx, query, userID, courseID).
		Scan(&res.ID, &res.CertificateCode, &res.UserID, &res.CourseID, &res.StudentName, &res.CourseTitle, &res.FinalScore, &res.IssuedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCertificateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}
