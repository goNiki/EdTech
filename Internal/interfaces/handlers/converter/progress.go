package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func LessonProgressToDTO(d *domain.LessonProgress) dto.LessonProgress {
	if d == nil {
		return dto.LessonProgress{}
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
	}
}

func CourseProgressToDTO(d *domain.CourseProgress) dto.CourseProgress {
	if d == nil {
		return dto.CourseProgress{}
	}
	return dto.CourseProgress{
		ID:             d.ID,
		CourseID:       d.CourseID,
		Status:         string(d.Status),
		Percent:        d.Percent,
		CompletedLess:  d.CompletedLess,
		StartedAt:      d.StartedAt,
		CompletedAt:    d.CompletedAt,
		LastAccessedAt: d.LastAccessedAt,
	}
}
