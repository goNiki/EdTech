package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func LessonProgressToDTO(d *domain.LessonProgress) dto.LessonProgress {
	if d == nil {
		return dto.LessonProgress{
			Submissions: make([]dto.LessonSubmissionDetailDTO, 0),
		}
	}
	submissions := make([]dto.LessonSubmissionDetailDTO, 0, len(d.Submissions))
	for _, s := range d.Submissions {
		submissions = append(submissions, dto.LessonSubmissionDetailDTO{
			QuestionText:    s.QuestionText,
			StudentAnswer:   s.StudentAnswer,
			PointsAwarded:   s.PointsAwarded,
			MaxPoints:       s.MaxPoints,
			TeacherFeedback: s.TeacherFeedback,
			IsGraded:        s.IsGraded,
		})
	}
	return dto.LessonProgress{
		ID:          d.ID,
		LessonID:    d.LessonID,
		Status:      string(d.Status),
		Score:       d.Score,
		TimeSpent:   d.TimeSpent,
		LastPos:     d.LastPos,
		StartedAt:   d.StartedAt,
		CompletedAt: d.CompletedAt,
		Submissions: submissions,
	}
}

func CourseProgressToDTO(d *domain.CourseProgress) dto.CourseProgress {
	if d == nil {
		return dto.CourseProgress{}
	}
	status := string(domain.ProgressStatusNotStarted)
	if d.CompletedAt != nil || (d.TotalLessons > 0 && d.CompletedLess >= d.TotalLessons) {
		status = string(domain.ProgressStatusCompleted)
	} else if d.StartedAt != nil || d.Percent > 0 || d.CompletedLess > 0 {
		status = string(domain.ProgressStatusInProgress)
	}

	return dto.CourseProgress{
		ID:             d.ID,
		CourseID:       d.CourseID,
		Status:         status,
		Percent:        d.Percent,
		CompletedLess:  d.CompletedLess,
		TotalLessons:   d.TotalLessons,
		AverageScore:   d.AverageScore,
		StartedAt:      d.StartedAt,
		CompletedAt:    d.CompletedAt,
		LastAccessedAt: d.LastAccessedAt,
	}
}
