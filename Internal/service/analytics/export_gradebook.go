package analytics

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"
)

func (s *analyticsService) ExportCourseGradebookCSV(ctx context.Context, teacherID, courseID int64) ([]byte, string, error) {
	const op = "service.analytics.ExportCourseGradebookCSV"

	if err := s.checkTeacherAccess(ctx, teacherID, courseID); err != nil {
		return nil, "", fmt.Errorf("%s: check access: %w", op, err)
	}

	records, err := s.analyticsRepo.GetCourseGradebook(ctx, courseID)
	if err != nil {
		return nil, "", fmt.Errorf("%s: get gradebook: %w", op, err)
	}

	var buf bytes.Buffer
	// Запись UTF-8 BOM для безупречного открытия в MS Excel
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	w := csv.NewWriter(&buf)
	w.Comma = ';'

	header := []string{
		"ID",
		"Студент",
		"Email",
		"Дата записи",
		"Прогресс (%)",
		"Пройдено уроков",
		"Всего уроков",
		"Средний балл (%)",
		"Статус",
		"Сертификат",
	}
	if err := w.Write(header); err != nil {
		return nil, "", fmt.Errorf("%s: write header: %w", op, err)
	}

	for _, rec := range records {
		certStr := "—"
		if rec.CertificateCode != nil && *rec.CertificateCode != "" {
			certStr = fmt.Sprintf("Выдан (%s)", *rec.CertificateCode)
		}

		row := []string{
			strconv.FormatInt(rec.UserID, 10),
			rec.StudentName,
			rec.Email,
			rec.EnrolledAt.Format("02.01.2006"),
			fmt.Sprintf("%d%%", rec.ProgressPercent),
			strconv.Itoa(rec.CompletedLessons),
			strconv.Itoa(rec.TotalLessons),
			fmt.Sprintf("%d%%", rec.AverageScore),
			rec.Status,
			certStr,
		}
		if err := w.Write(row); err != nil {
			return nil, "", fmt.Errorf("%s: write row: %w", op, err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", fmt.Errorf("%s: flush csv: %w", op, err)
	}

	filename := fmt.Sprintf("gradebook_course_%d_%s.csv", courseID, time.Now().Format("2006-01-02"))

	return buf.Bytes(), filename, nil
}
