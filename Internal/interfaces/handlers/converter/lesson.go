package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func LessonToDTO(d *domain.Lesson) dto.Lesson {
	if d == nil {
		return dto.Lesson{}
	}
	return dto.Lesson{
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
		QuizSettings: d.QuizSettings,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		PublishedAt:  d.PublishedAt,
	}
}

func CreateLessonRequestToDomain(req dto.CreateLessonRequest) domain.Lesson {
	status := req.Status
	if status == "" {
		status = domain.StatusDraft
	}
	return domain.Lesson{
		CourseID:     req.CourseID,
		SectionID:    req.SectionID,
		Title:        req.Title,
		Description:  req.Description,
		CoverURL:     req.CoverURL,
		Content:      req.Content,
		Type:         req.Type,
		Duration:     req.Duration,
		IsFree:       req.IsFree,
		Status:       status,
		QuizSettings: req.QuizSettings,
	}
}

func LessonNavigationToDTO(nav *domain.LessonNavigationContext) dto.LessonNavigationResponse {
	if nav == nil {
		return dto.LessonNavigationResponse{
			Code:    200,
			Message: "success",
		}
	}

	var prevDTO *dto.LessonNavNeighborDTO
	if nav.PrevLesson != nil {
		prevDTO = &dto.LessonNavNeighborDTO{
			ID:    nav.PrevLesson.ID,
			Title: nav.PrevLesson.Title,
		}
	}

	var nextDTO *dto.LessonNavNeighborDTO
	if nav.NextLesson != nil {
		nextDTO = &dto.LessonNavNeighborDTO{
			ID:    nav.NextLesson.ID,
			Title: nav.NextLesson.Title,
		}
	}

	syllabusDTO := make([]dto.LessonNavSectionDTO, 0, len(nav.Syllabus))
	for _, sec := range nav.Syllabus {
		lessonsDTO := make([]dto.LessonNavItemDTO, 0, len(sec.Lessons))
		for _, l := range sec.Lessons {
			lessonsDTO = append(lessonsDTO, dto.LessonNavItemDTO{
				ID:          l.ID,
				Title:       l.Title,
				Position:    l.Position,
				IsCompleted: l.IsCompleted,
				Score:       l.Score,
			})
		}
		syllabusDTO = append(syllabusDTO, dto.LessonNavSectionDTO{
			SectionID:    sec.SectionID,
			SectionTitle: sec.SectionTitle,
			Position:     sec.Position,
			Lessons:      lessonsDTO,
		})
	}

	return dto.LessonNavigationResponse{
		Code:    200,
		Message: "success",
		Data: dto.LessonNavigationDataDTO{
			CurrentLesson: dto.LessonNavCurrentDTO{
				ID:        nav.CurrentLesson.ID,
				Title:     nav.CurrentLesson.Title,
				Position:  nav.CurrentLesson.Position,
				SectionID: nav.CurrentLesson.SectionID,
			},
			Course: dto.LessonNavCourseDTO{
				ID:    nav.Course.ID,
				Title: nav.Course.Title,
				Slug:  nav.Course.Slug,
			},
			PrevLesson: prevDTO,
			NextLesson: nextDTO,
			Syllabus:   syllabusDTO,
		},
	}
}

