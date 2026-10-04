package course

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

const baseCourseDetailSelect = `
	SELECT 
		c.id, 
		c.title, 
		c.slug, 
		c.short_description, 
		c.description, 
		c.cover_url, 
		c.intro_video_url, 
		c.created_by, 
		c.visibility, 
		c.status, 
		c.difficulty, 
		c.language, 
		c.estimated_duration, 
		c.category_id, 
		c.total_lessons, 
		c.total_sections, 
		c.enrolled_count, 
		c.created_at, 
		c.updated_at, 
		c.published_at, 
		c.archived_at, 
		c.deleted_at,
		u.id AS author_id,
		COALESCE(NULLIF(TRIM(u.first_name || ' ' || u.last_name), ''), u.username) AS author_name,
		u.avatar_url AS author_avatar_url,
		COALESCE(u.headline, '') AS author_headline,
		COALESCE(u.bio, '') AS author_bio,
		(SELECT COUNT(*) FROM courses WHERE created_by = u.id AND status = 'published' AND deleted_at IS NULL) AS author_courses_count,
		(SELECT COALESCE(SUM(enrolled_count), 0) FROM courses WHERE created_by = u.id AND deleted_at IS NULL) AS author_total_students
	FROM courses c
	LEFT JOIN users u ON c.created_by = u.id
`

func scanCourseWithAuthor(row pgx.Row) (*domain.Course, error) {
	var (
		entity              repomodels.Course
		authorID            *int64
		authorName          *string
		authorAvatarURL     *string
		authorHeadline      string
		authorBio           string
		authorCoursesCount  int
		authorTotalStudents int
	)

	err := row.Scan(
		&entity.ID,
		&entity.Title,
		&entity.Slug,
		&entity.ShortDescription,
		&entity.Description,
		&entity.CoverURL,
		&entity.IntroVideoURL,
		&entity.CreatedBy,
		&entity.Visibility,
		&entity.Status,
		&entity.Difficulty,
		&entity.Language,
		&entity.EstimatedDuration,
		&entity.CategoryID,
		&entity.TotalLessons,
		&entity.TotalSections,
		&entity.EnrolledCount,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.PublishedAt,
		&entity.ArchivedAt,
		&entity.DeletedAt,
		&authorID,
		&authorName,
		&authorAvatarURL,
		&authorHeadline,
		&authorBio,
		&authorCoursesCount,
		&authorTotalStudents,
	)
	if err != nil {
		return nil, err
	}

	domainCourse := repoconverter.CourseToDomain(&entity)
	if authorID != nil && *authorID > 0 {
		name := ""
		if authorName != nil {
			name = *authorName
		}
		domainCourse.Author = &domain.CourseAuthorInfo{
			ID:            *authorID,
			Name:          name,
			AvatarURL:     authorAvatarURL,
			Headline:      authorHeadline,
			Bio:           authorBio,
			CoursesCount:  authorCoursesCount,
			TotalStudents: authorTotalStudents,
		}
	}

	return domainCourse, nil
}

func (r *repository) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	const op = "repository.course.GetCourseByID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := baseCourseDetailSelect + ` WHERE c.id = $1 AND c.deleted_at IS NULL`

	row := q.QueryRow(ctx, query, id)
	course, err := scanCourseWithAuthor(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return course, nil
}

func (r *repository) GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error) {
	const op = "repository.course.GetCourseBySlug"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := baseCourseDetailSelect + ` WHERE c.slug = $1 AND c.deleted_at IS NULL`

	row := q.QueryRow(ctx, query, slug)
	course, err := scanCourseWithAuthor(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return course, nil
}
