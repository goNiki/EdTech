package quiz_test

import (
	"strings"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/service/quiz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeLessonContentForStudent(t *testing.T) {
	rawJSON := `{
		"content": [
			{
				"type": "QuizSingleBlock",
				"props": {
					"id": "single-1",
					"question": "What is Go?",
					"options": [
						{"text": "A compiled language", "isCorrect": true, "explain": "Because it compiles to machine code"},
						{"text": "An interpreted script", "isCorrect": false}
					]
				}
			},
			{
				"type": "QuizMatchBlock",
				"props": {
					"id": "match-1",
					"pairs": [
						{"left": "Zebra", "right": "Animal"},
						{"left": "Apple", "right": "Fruit"}
					]
				}
			},
			{
				"type": "QuizDropdownBlankBlock",
				"props": {
					"id": "drop-1",
					"templateText": "We use {Postgres; Redis, Mongo} for SQL",
					"correctIndex": 0,
					"blanks": [
						{"key": "blank_1", "correctAnswer": "Postgres"}
					]
				}
			},
			{
				"type": "QuizInputBlankBlock",
				"props": {
					"id": "input-1",
					"templateText": "Write {secret_answer} here",
					"correctAnswer": "secret_answer"
				}
			},
			{
				"type": "QuizSequenceBlock",
				"props": {
					"id": "seq-1",
					"items": [
						{"text": "Step 1"},
						{"text": "Step 2"},
						{"text": "Step 3"}
					]
				}
			}
		]
	}`

	sanitized, err := quiz.SanitizeLessonContentForStudent(rawJSON)
	require.NoError(t, err)

	// Invariants check
	assert.NotContains(t, sanitized, "isCorrect")
	assert.NotContains(t, sanitized, "is_correct")
	assert.NotContains(t, sanitized, "explain")
	assert.NotContains(t, sanitized, "correctIndex")
	assert.NotContains(t, sanitized, "secret_answer")
	assert.NotContains(t, sanitized, "correctAnswer")

	// Input blank replacement
	assert.Contains(t, sanitized, "{blank_1}")

	// Dropdown options reordering
	assert.Contains(t, sanitized, "Mongo")
	assert.Contains(t, sanitized, "Postgres")
	assert.Contains(t, sanitized, "Redis")
	// "Mongo" comes before "Postgres" alphabetically, so "{Postgres;" is gone
	assert.False(t, strings.Contains(sanitized, "{Postgres;"))
}

func TestValidateQuizSubmission(t *testing.T) {
	rawJSON := `{
		"content": [
			{
				"type": "QuizSingleBlock",
				"props": {
					"id": "single-1",
					"points": 10,
					"options": [
						{"text": "Wrong", "isCorrect": false},
						{"text": "Correct", "isCorrect": true}
					]
				}
			},
			{
				"type": "QuizMultiBlock",
				"props": {
					"id": "multi-1",
					"points": 20,
					"options": [
						{"text": "Correct 1", "isCorrect": true},
						{"text": "Wrong", "isCorrect": false},
						{"text": "Correct 2", "isCorrect": true}
					]
				}
			},
			{
				"type": "QuizInputBlankBlock",
				"props": {
					"id": "input-1",
					"points": 10,
					"templateText": "Language is {Go}"
				}
			}
		]
	}`

	t.Run("All correct answers", func(t *testing.T) {
		answers := []domain.LessonAnswerSubmission{
			{
				BlockID: "single-1",
				Answer:  map[string]any{"selected_option": 1},
			},
			{
				BlockID: "multi-1",
				Answer:  map[string]any{"selected_options": []any{0, 2}},
			},
			{
				BlockID: "input-1",
				Answer:  map[string]any{"blanks": map[string]any{"blank_1": "Go"}},
			},
		}

		res, err := quiz.ValidateQuizSubmission(rawJSON, answers)
		require.NoError(t, err)
		assert.Equal(t, 100, res.Score)
		assert.Equal(t, 40, res.EarnedPoints)
		assert.Equal(t, 40, res.TotalMaxPoints)
		assert.True(t, res.Results["single-1"].IsCorrect)
		assert.True(t, res.Results["multi-1"].IsCorrect)
		assert.True(t, res.Results["input-1"].IsCorrect)
	})

	t.Run("Partial correct answers", func(t *testing.T) {
		answers := []domain.LessonAnswerSubmission{
			{
				BlockID: "single-1",
				Answer:  map[string]any{"selected_option": 1}, // 10 pts
			},
			{
				BlockID: "multi-1",
				Answer:  map[string]any{"selected_options": []any{0}}, // wrong (20 pts lost)
			},
			{
				BlockID: "input-1",
				Answer:  map[string]any{"blanks": map[string]any{"blank_1": "Wrong"}}, // wrong (10 pts lost)
			},
		}

		res, err := quiz.ValidateQuizSubmission(rawJSON, answers)
		require.NoError(t, err)
		// 10 out of 40 = 25%
		assert.Equal(t, 25, res.Score)
		assert.Equal(t, 10, res.EarnedPoints)
		assert.Equal(t, 40, res.TotalMaxPoints)
		assert.True(t, res.Results["single-1"].IsCorrect)
		assert.False(t, res.Results["multi-1"].IsCorrect)
		assert.False(t, res.Results["input-1"].IsCorrect)
	})
}
