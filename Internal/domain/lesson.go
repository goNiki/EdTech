package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	errorsAPP "edtech/pkg/errors"
)

type Lesson struct {
	ID          int64
	CourseID    int64
	SectionID   *int64
	Title       string
	Description string
	CoverURL    string
	Content     string
	Type        string
	Position    int64
	Duration    *int
	IsFree      bool
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
	DeletedAt   *time.Time
}

func ValidateLessonContent(content string) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil
	}

	if !json.Valid([]byte(trimmed)) {
		return fmt.Errorf("%w: invalid json structure", errorsAPP.ErrLessonValidation)
	}

	var root interface{}
	if err := json.Unmarshal([]byte(trimmed), &root); err != nil {
		return fmt.Errorf("%w: unmarshal content json: %v", errorsAPP.ErrLessonValidation, err)
	}

	var rawBlocks []interface{}
	switch v := root.(type) {
	case map[string]interface{}:
		if contentVal, ok := v["content"]; ok {
			if slice, ok := contentVal.([]interface{}); ok {
				rawBlocks = slice
			} else if contentVal != nil {
				return fmt.Errorf("%w: 'content' field must be an array of blocks", errorsAPP.ErrLessonValidation)
			}
		}
	case []interface{}:
		rawBlocks = v
	default:
		return fmt.Errorf("%w: root content must be a JSON object or array", errorsAPP.ErrLessonValidation)
	}

	seenIDs := make(map[string]bool, len(rawBlocks))
	for i, b := range rawBlocks {
		block, ok := b.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: block at index %d is not a JSON object", errorsAPP.ErrLessonValidation, i)
		}

		blockType, ok := block["type"].(string)
		if !ok || strings.TrimSpace(blockType) == "" {
			return fmt.Errorf("%w: block at index %d missing required field 'type'", errorsAPP.ErrLessonValidation, i)
		}

		var blockID string
		if props, ok := block["props"].(map[string]interface{}); ok {
			if idVal, ok := props["id"].(string); ok {
				blockID = strings.TrimSpace(idVal)
			}
		}
		if blockID == "" {
			if idVal, ok := block["id"].(string); ok {
				blockID = strings.TrimSpace(idVal)
			}
		}

		lowerType := strings.ToLower(blockType)
		isQuizBlock := strings.Contains(lowerType, "quiz") ||
			strings.Contains(lowerType, "choice") ||
			strings.Contains(lowerType, "matching") ||
			strings.Contains(lowerType, "question")

		if isQuizBlock && blockID == "" {
			return fmt.Errorf("%w: quiz block %q at index %d must have non-empty props.id", errorsAPP.ErrLessonValidation, blockType, i)
		}

		if blockID != "" {
			if seenIDs[blockID] {
				return fmt.Errorf("%w: duplicate block id %q found at index %d", errorsAPP.ErrLessonValidation, blockID, i)
			}
			seenIDs[blockID] = true
		}
	}

	return nil
}

func (l *Lesson) Validate() error {
	if l.CourseID == 0 || l.Title == "" || l.Description == "" {
		return errorsAPP.ErrLessonValidation
	}
	if l.Content != "" {
		if err := ValidateLessonContent(l.Content); err != nil {
			return err
		}
	}
	return nil
}

func (l *Lesson) CanPublish() error {
	if l.Status == StatusPublished {
		return errorsAPP.ErrLessonAlreadyPublished
	}
	return nil
}

func (l *Lesson) Publish(now time.Time) {
	l.Status = StatusPublished
	l.PublishedAt = &now
	l.UpdatedAt = now
}

func (l *Lesson) CanArchive() error {
	if l.Status == StatusArchived {
		return errorsAPP.ErrLessonAlreadyArchived
	}
	return nil
}

func (l *Lesson) Archive(now time.Time) {
	l.Status = StatusArchived
	l.UpdatedAt = now
}

func (l *Lesson) IsPublished() bool {
	return l.Status == StatusPublished
}
