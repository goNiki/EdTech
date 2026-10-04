package certificate

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"edtech/internal/domain"
	"edtech/internal/repository"
	services "edtech/internal/service"
	errorsAPP "edtech/pkg/errors"
)

type service struct {
	certRepo     repository.CertificateRepository
	courseRepo   repository.CourseRepository
	userRepo     repository.UserRepository
	progressRepo repository.ProgressRepository
}

func NewCertificateService(
	certRepo repository.CertificateRepository,
	courseRepo repository.CourseRepository,
	userRepo repository.UserRepository,
	progressRepo repository.ProgressRepository,
) services.CertificateServices {
	return &service{
		certRepo:     certRepo,
		courseRepo:   courseRepo,
		userRepo:     userRepo,
		progressRepo: progressRepo,
	}
}

func (s *service) GetOrIssueCertificate(ctx context.Context, userID, courseID int64) (*domain.Certificate, error) {
	const op = "service.certificate.GetOrIssueCertificate"

	// 1. Check if certificate already exists (idempotency)
	existing, err := s.certRepo.GetCertificateByUserAndCourse(ctx, userID, courseID)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && !errors.Is(err, errorsAPP.ErrCertificateNotFound) {
		return nil, fmt.Errorf("%s: check existing: %w", op, err)
	}

	// 2. Check course progress - must be 100% completed
	progress, err := s.progressRepo.GetCourseProgress(ctx, userID, courseID)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrCourseProgressNotFound) || errors.Is(err, errorsAPP.ErrProgressNotFound) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotCompleted)
		}
		return nil, fmt.Errorf("%s: get course progress: %w", op, err)
	}

	if progress.Percent < 100 || (progress.TotalLessons > 0 && progress.CompletedLess < progress.TotalLessons) {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotCompleted)
	}

	// 3. Fetch user details for certificate name
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: get user: %w", op, err)
	}

	studentName := formatStudentName(user)

	// 4. Fetch course details
	course, err := s.courseRepo.GetCourseByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: get course: %w", op, err)
	}

	// 5. Generate unique certificate code: EDL-YYYY-XXXXXXXX
	code, err := generateCertificateCode()
	if err != nil {
		return nil, fmt.Errorf("%s: generate code: %w", op, err)
	}

	finalScore := 0.0
	if progress.AverageScore != nil {
		finalScore = *progress.AverageScore
	}

	cert := &domain.Certificate{
		CertificateCode: code,
		UserID:          userID,
		CourseID:        courseID,
		StudentName:     studentName,
		CourseTitle:     course.Title,
		FinalScore:      finalScore,
	}

	created, err := s.certRepo.CreateCertificate(ctx, cert)
	if err != nil {
		return nil, fmt.Errorf("%s: save certificate: %w", op, err)
	}

	return created, nil
}

func (s *service) VerifyCertificate(ctx context.Context, code string) (*domain.Certificate, error) {
	const op = "service.certificate.VerifyCertificate"

	trimmed := strings.TrimSpace(code)
	if trimmed == "" {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCertificateNotFound)
	}

	cert, err := s.certRepo.GetCertificateByCode(ctx, trimmed)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cert, nil
}

func formatStudentName(user *domain.User) string {
	var parts []string
	if user.FirstName != nil && *user.FirstName != "" {
		parts = append(parts, *user.FirstName)
	}
	if user.LastName != nil && *user.LastName != "" {
		parts = append(parts, *user.LastName)
	}
	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}
	return user.Username
}

func generateCertificateCode() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	year := time.Now().Year()
	return fmt.Sprintf("EDL-%d-%s", year, strings.ToUpper(hex.EncodeToString(bytes))), nil
}
