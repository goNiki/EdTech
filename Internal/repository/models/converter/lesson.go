package converter

import (
	"encoding/json"

	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func LessonToDomain(m *models.Lesson) *domain.Lesson {
	if m == nil {
		return nil
	}
	var qs *domain.QuizSettings
	if len(m.QuizSettings) > 0 {
		var s domain.QuizSettings
		if err := json.Unmarshal(m.QuizSettings, &s); err == nil {
			qs = &s
		}
	}
	return &domain.Lesson{
		ID:           m.ID,
		CourseID:     m.CourseID,
		SectionID:    m.SectionID,
		Title:        m.Title,
		Description:  m.Description,
		CoverURL:     m.CoverURL,
		Content:      m.Content,
		Type:         m.Type,
		Position:     m.Position,
		Duration:     m.Duration,
		IsFree:       m.IsFree,
		Status:       m.Status,
		QuizSettings: qs,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		PublishedAt:  m.PublishedAt,
		DeletedAt:    m.DeletedAt,
	}
}

func LessonToEntity(d *domain.Lesson) *models.Lesson {
	if d == nil {
		return nil
	}
	var qsBytes []byte
	if d.QuizSettings != nil {
		qsBytes, _ = json.Marshal(d.QuizSettings)
	}
	return &models.Lesson{
		ID:           d.ID,
		CourseID:     d.CourseID,
		SectionID:    d.SectionID,
		Title:        d.Title,
		Description:  d.Description,
		CoverURL:     d.CoverURL,
		Content:      d.Content,
		Type:         d.Type,
		Position:     d.Position,
		Duration:     d.Duration,
		IsFree:       d.IsFree,
		Status:       d.Status,
		QuizSettings: qsBytes,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		PublishedAt:  d.PublishedAt,
		DeletedAt:    d.DeletedAt,
	}
}
