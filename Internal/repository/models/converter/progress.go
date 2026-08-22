package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func LessonProgressToDomain(m *models.LessonProgress) *domain.LessonProgress {
	if m == nil { return nil }
	return &domain.LessonProgress{
		ID:          m.ID,
		UserID:      m.UserID,
		LessonID:    m.LessonID,
		Status:      domain.ProgressStatus(m.Status),
		Score:       m.Score,
		TimeSpent:   m.TimeSpent,
		LastPos:     m.LastPos,
		StartedAt:   m.StartedAt,
		CompletedAt: m.CompletedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func LessonProgressToEntity(d *domain.LessonProgress) *models.LessonProgress {
	if d == nil { return nil }
	return &models.LessonProgress{
		ID:          d.ID,
		UserID:      d.UserID,
		LessonID:    d.LessonID,
		Status:      string(d.Status),
		Score:       d.Score,
		TimeSpent:   d.TimeSpent,
		LastPos:     d.LastPos,
		StartedAt:   d.StartedAt,
		CompletedAt: d.CompletedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func CourseProgressToDomain(m *models.CourseProgress) *domain.CourseProgress {
	if m == nil { return nil }
	return &domain.CourseProgress{
		ID:             m.ID,
		UserID:         m.UserID,
		CourseID:       m.CourseID,
		Status:         domain.ProgressStatus(m.Status),
		Percent:        m.Percent,
		CompletedLess:  m.CompletedLess,
		StartedAt:      m.StartedAt,
		CompletedAt:    m.CompletedAt,
		LastAccessedAt: m.LastAccessedAt,
	}
}

func CourseProgressToEntity(d *domain.CourseProgress) *models.CourseProgress {
	if d == nil { return nil }
	return &models.CourseProgress{
		ID:             d.ID,
		UserID:         d.UserID,
		CourseID:       d.CourseID,
		Status:         string(d.Status),
		Percent:        d.Percent,
		CompletedLess:  d.CompletedLess,
		StartedAt:      d.StartedAt,
		CompletedAt:    d.CompletedAt,
		LastAccessedAt: d.LastAccessedAt,
	}
}
