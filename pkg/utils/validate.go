package utils

import (
	"edtech/internal/dto"
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"strings"
)

// TODO сделать норм валидацию.
func ValidateCourse(title string, slug string, createdBy int64, visibility string, status string) error {

	var fields []string

	if title == "" {
		fields = append(fields, "title")
	}

	if slug == "" {
		fields = append(fields, "slug")
	}

	if createdBy == 0 {
		fields = append(fields, "createdBy")
	}

	if visibility != "private" && visibility != "public" {
		fields = append(fields, "visibility")
	}

	if status != "draft" && status != "published" {
		fields = append(fields, "status")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%w, %s are required/invalid", errorsAPP.ErrCourseValidation, strings.Join(fields, ", "))
	}

	return nil
}

func ValidateLesson(courseid int64, title string, description string) error {

	var fields []string

	if courseid == 0 {
		fields = append(fields, "courseid")
	}

	if title == "" {
		fields = append(fields, "title")
	}

	if description == "" {
		fields = append(fields, "description")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%w: %s are required/invalid", errorsAPP.ErrLessonValidation, strings.Join(fields, ", "))
	}
	return nil
}

func ValidateEnrolle(enrol dto.EnrollRequest) error {
	var fields []string

	if enrol.UserEmail == "" {
		fields = append(fields, "userEmail")
	}

	if enrol.CourseID == 0 {
		fields = append(fields, "courseID")
	}

	if enrol.Role != "student" && enrol.Role != "teacher" {
		fields = append(fields, "role")
	}

	if len(fields) > 0 {
		return fmt.Errorf("%w: %s", errorsAPP.ErrFailEnroleValidate, strings.Join(fields, ", "))
	}

	return nil
}
