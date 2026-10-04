package quiz

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"edtech/internal/domain"
)

var (
	// dropdownRegex matches {Correct; Distractor1, Distractor2} or {Word}
	dropdownRegex = regexp.MustCompile(`\{([^{}]+)\}`)
)

// SanitizeLessonContentForStudent removes all answer keys and distractor order hints from Puck content JSON.
func SanitizeLessonContentForStudent(contentJSON string) (string, error) {
	if strings.TrimSpace(contentJSON) == "" {
		return contentJSON, nil
	}

	var root any
	if err := json.Unmarshal([]byte(contentJSON), &root); err != nil {
		// If not valid JSON, return as-is
		return contentJSON, nil
	}

	switch v := root.(type) {
	case map[string]any:
		if contentSlice, ok := v["content"].([]any); ok {
			for _, item := range contentSlice {
				if blockMap, ok := item.(map[string]any); ok {
					sanitizeBlock(blockMap)
				}
			}
		}
	case []any:
		for _, item := range v {
			if blockMap, ok := item.(map[string]any); ok {
				sanitizeBlock(blockMap)
			}
		}
	}

	sanitizedBytes, err := json.Marshal(root)
	if err != nil {
		return "", fmt.Errorf("marshal sanitized content: %w", err)
	}

	return string(sanitizedBytes), nil
}

func sanitizeBlock(block map[string]any) {
	blockType, _ := block["type"].(string)
	props, ok := block["props"].(map[string]any)
	if !ok || props == nil {
		return
	}

	switch blockType {
	case "QuizSingleBlock", "QuizMultiBlock":
		sanitizeChoiceBlock(props)
	case "QuizMatchBlock":
		sanitizeMatchBlock(props)
	case "QuizDropdownBlankBlock":
		sanitizeDropdownBlock(props)
	case "QuizInputBlankBlock":
		sanitizeInputBlankBlock(props)
	case "QuizSequenceBlock":
		sanitizeSequenceBlock(props)
	}

	recursiveCleanSensitiveKeys(props)
}

func sanitizeChoiceBlock(props map[string]any) {
	opts, ok := props["options"].([]any)
	if !ok {
		return
	}
	for _, opt := range opts {
		if optMap, ok := opt.(map[string]any); ok {
			delete(optMap, "isCorrect")
			delete(optMap, "is_correct")
			delete(optMap, "correct")
			delete(optMap, "explain")
		}
	}
}

func sanitizeMatchBlock(props map[string]any) {
	pairs, ok := props["pairs"].([]any)
	if !ok || len(pairs) == 0 {
		return
	}
	rightValues := make([]string, 0, len(pairs))
	for _, p := range pairs {
		if pMap, ok := p.(map[string]any); ok {
			delete(pMap, "correctPair")
			delete(pMap, "correct_pair")
			delete(pMap, "correct")
			if r, ok := pMap["right"].(string); ok {
				rightValues = append(rightValues, r)
			}
		}
	}
	sort.Strings(rightValues)
	for i, p := range pairs {
		if pMap, ok := p.(map[string]any); ok && i < len(rightValues) {
			pMap["right"] = rightValues[i]
		}
	}
}

func sanitizeDropdownBlock(props map[string]any) {
	if tmpl, ok := props["templateText"].(string); ok && tmpl != "" {
		props["templateText"] = sanitizeDropdownTemplate(tmpl)
	}
	delete(props, "correctIndex")
	delete(props, "correct_index")
	delete(props, "correct")
	if blanks, ok := props["blanks"].([]any); ok {
		for _, blk := range blanks {
			if bMap, ok := blk.(map[string]any); ok {
				delete(bMap, "correctAnswer")
				delete(bMap, "correct_answer")
				delete(bMap, "correct")
			}
		}
	}
}

func sanitizeInputBlankBlock(props map[string]any) {
	if tmpl, ok := props["templateText"].(string); ok && tmpl != "" {
		counter := 1
		sanitizedTmpl := dropdownRegex.ReplaceAllStringFunc(tmpl, func(_ string) string {
			tag := fmt.Sprintf("{blank_%d}", counter)
			counter++
			return tag
		})
		props["templateText"] = sanitizedTmpl
	}
	delete(props, "correctAnswer")
	delete(props, "correct_answer")
	delete(props, "correct")
	if blanks, ok := props["blanks"].([]any); ok {
		for _, blk := range blanks {
			if bMap, ok := blk.(map[string]any); ok {
				delete(bMap, "correctAnswer")
				delete(bMap, "correct_answer")
				delete(bMap, "correct")
			}
		}
	}
}

func sanitizeSequenceBlock(props map[string]any) {
	items, ok := props["items"].([]any)
	if !ok || len(items) <= 1 {
		return
	}
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	props["items"] = items
}

func sanitizeDropdownTemplate(template string) string {
	return dropdownRegex.ReplaceAllStringFunc(template, func(tag string) string {
		inner := strings.Trim(tag, "{}")
		if !strings.Contains(inner, ";") {
			return tag
		}

		parts := strings.Split(inner, ";")
		first := strings.TrimSpace(parts[0])
		var options []string
		if len(parts) > 1 {
			for _, opt := range strings.Split(parts[1], ",") {
				trimmed := strings.TrimSpace(opt)
				if trimmed != "" {
					options = append(options, trimmed)
				}
			}
		}
		if first != "" {
			options = append(options, first)
		}

		uniqueOptions := deduplicateStrings(options)
		sort.Strings(uniqueOptions)

		if len(uniqueOptions) > 1 {
			return fmt.Sprintf("{%s; %s}", uniqueOptions[0], strings.Join(uniqueOptions[1:], ", "))
		}
		if len(uniqueOptions) == 1 {
			return fmt.Sprintf("{%s}", uniqueOptions[0])
		}
		return tag
	})
}

func deduplicateStrings(items []string) []string {
	seen := make(map[string]bool)
	var res []string
	for _, it := range items {
		if !seen[it] {
			seen[it] = true
			res = append(res, it)
		}
	}
	return res
}

func recursiveCleanSensitiveKeys(m map[string]any) {
	sensitive := map[string]bool{
		"isCorrect":      true,
		"is_correct":     true,
		"correct":        true,
		"correctAnswer":  true,
		"correct_answer": true,
		"correctPair":    true,
		"correct_pair":   true,
		"correctIndex":   true,
		"correct_index":  true,
	}

	for k, v := range m {
		if sensitive[k] {
			delete(m, k)
			continue
		}
		if subMap, ok := v.(map[string]any); ok {
			recursiveCleanSensitiveKeys(subMap)
		} else if slice, ok := v.([]any); ok {
			for _, item := range slice {
				if itemMap, ok := item.(map[string]any); ok {
					recursiveCleanSensitiveKeys(itemMap)
				}
			}
		}
	}
}

// ValidateQuizSubmission calculates the actual score and per-block results on the server based on canonical lesson content.
func ValidateQuizSubmission(contentJSON string, answers []domain.LessonAnswerSubmission) (*domain.LessonCompletionResult, error) {
	answersMap := make(map[string]any)
	for _, a := range answers {
		answersMap[a.BlockID] = a.Answer
	}

	var root any
	if err := json.Unmarshal([]byte(contentJSON), &root); err != nil {
		// Non-JSON content has no quizzes to validate
		return &domain.LessonCompletionResult{
			Score:          100,
			EarnedPoints:   0,
			TotalMaxPoints: 0,
			Results:        make(map[string]domain.BlockValidationResult),
		}, nil
	}

	var blocks []map[string]any
	switch v := root.(type) {
	case map[string]any:
		if contentSlice, ok := v["content"].([]any); ok {
			for _, item := range contentSlice {
				if blockMap, ok := item.(map[string]any); ok {
					blocks = append(blocks, blockMap)
				}
			}
		}
	case []any:
		for _, item := range v {
			if blockMap, ok := item.(map[string]any); ok {
				blocks = append(blocks, blockMap)
			}
		}
	}

	totalMaxPoints := 0
	earnedPoints := 0
	results := make(map[string]domain.BlockValidationResult)

	for idx, block := range blocks {
		blockType, _ := block["type"].(string)
		props, ok := block["props"].(map[string]any)
		if !ok || props == nil {
			continue
		}

		blockID, _ := props["id"].(string)
		if blockID == "" {
			blockID = fmt.Sprintf("block-%d", idx)
		}

		points := parsePoints(props["points"], defaultPointsForType(blockType))

		switch blockType {
		case "QuizSingleBlock":
			totalMaxPoints += points
			res := validateSingleBlock(props, answersMap[blockID], points)
			results[blockID] = res
			if res.IsCorrect {
				earnedPoints += points
			}

		case "QuizMultiBlock":
			totalMaxPoints += points
			res := validateMultiBlock(props, answersMap[blockID], points)
			results[blockID] = res
			if res.IsCorrect {
				earnedPoints += points
			}

		case "QuizMatchBlock":
			totalMaxPoints += points
			res, pts := validateMatchBlock(props, answersMap[blockID], points)
			results[blockID] = res
			earnedPoints += pts

		case "QuizDropdownBlankBlock":
			totalMaxPoints += points
			res, pts := validateDropdownBlock(props, answersMap[blockID], points)
			results[blockID] = res
			earnedPoints += pts

		case "QuizInputBlankBlock":
			totalMaxPoints += points
			res, pts := validateInputBlankBlock(props, answersMap[blockID], points)
			results[blockID] = res
			earnedPoints += pts

		case "QuizSequenceBlock":
			totalMaxPoints += points
			res := validateSequenceBlock(props, answersMap[blockID])
			results[blockID] = res
			if res.IsCorrect {
				earnedPoints += points
			}
		}
	}

	finalScore := 100
	if totalMaxPoints > 0 {
		finalScore = int(math.Round((float64(earnedPoints) / float64(totalMaxPoints)) * 100))
	}

	return &domain.LessonCompletionResult{
		Score:          finalScore,
		EarnedPoints:   earnedPoints,
		TotalMaxPoints: totalMaxPoints,
		Results:        results,
	}, nil
}

func defaultPointsForType(blockType string) int {
	switch blockType {
	case "QuizSingleBlock":
		return 10
	case "QuizMultiBlock":
		return 15
	case "QuizMatchBlock":
		return 20
	case "QuizDropdownBlankBlock":
		return 10
	case "QuizInputBlankBlock":
		return 10
	case "QuizSequenceBlock":
		return 15
	default:
		return 10
	}
}

func parsePoints(val any, defaultVal int) int {
	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		if p, err := strconv.Atoi(v); err == nil {
			return p
		}
	}
	return defaultVal
}

func validateSingleBlock(props map[string]any, studentAnswer any, _ int) domain.BlockValidationResult {
	correctIdx := -1
	if opts, ok := props["options"].([]any); ok {
		for i, opt := range opts {
			if optMap, ok := opt.(map[string]any); ok {
				if isBoolCorrect(optMap["isCorrect"]) || isBoolCorrect(optMap["is_correct"]) {
					correctIdx = i
					break
				}
			}
		}
	}

	selected := extractIntOption(studentAnswer)
	if selected != nil && *selected == correctIdx && correctIdx != -1 {
		return domain.BlockValidationResult{
			IsCorrect:     true,
			Feedback:      "Верный ответ",
			CorrectAnswer: correctIdx,
		}
	}

	return domain.BlockValidationResult{
		IsCorrect:     false,
		Feedback:      "Неверный ответ",
		CorrectAnswer: correctIdx,
	}
}

func validateMultiBlock(props map[string]any, studentAnswer any, _ int) domain.BlockValidationResult {
	var correctIndices []int
	if opts, ok := props["options"].([]any); ok {
		for i, opt := range opts {
			if optMap, ok := opt.(map[string]any); ok {
				if isBoolCorrect(optMap["isCorrect"]) || isBoolCorrect(optMap["is_correct"]) {
					correctIndices = append(correctIndices, i)
				}
			}
		}
	}

	selected := extractIntSlice(studentAnswer)
	sort.Ints(selected)
	sort.Ints(correctIndices)

	isEqual := len(selected) == len(correctIndices)
	if isEqual {
		for i := range selected {
			if selected[i] != correctIndices[i] {
				isEqual = false
				break
			}
		}
	}

	return domain.BlockValidationResult{
		IsCorrect:     isEqual,
		Feedback:      formatBoolFeedback(isEqual),
		CorrectAnswer: correctIndices,
	}
}

func validateMatchBlock(props map[string]any, studentAnswer any, maxPoints int) (domain.BlockValidationResult, int) {
	pairs, ok := props["pairs"].([]any)
	if !ok || len(pairs) == 0 {
		return domain.BlockValidationResult{IsCorrect: true}, maxPoints
	}

	userPairs := extractStringMap(studentAnswer)
	correctCount := 0

	for i, p := range pairs {
		if pMap, ok := p.(map[string]any); ok {
			expected, _ := pMap["right"].(string)
			actual := userPairs[strconv.Itoa(i)]
			if strings.TrimSpace(strings.ToLower(actual)) == strings.TrimSpace(strings.ToLower(expected)) {
				correctCount++
			}
		}
	}

	awarded := int(math.Round(float64(correctCount) / float64(len(pairs)) * float64(maxPoints)))
	isCorrect := (correctCount == len(pairs))

	return domain.BlockValidationResult{
		IsCorrect: isCorrect,
		Feedback:  fmt.Sprintf("Верно %d из %d", correctCount, len(pairs)),
	}, awarded
}

func validateDropdownBlock(props map[string]any, studentAnswer any, maxPoints int) (domain.BlockValidationResult, int) {
	template, _ := props["templateText"].(string)
	matches := dropdownRegex.FindAllStringSubmatch(template, -1)
	if len(matches) == 0 {
		return domain.BlockValidationResult{IsCorrect: true}, maxPoints
	}

	userBlanks := extractStringMap(studentAnswer)
	correctCount := 0

	for i, match := range matches {
		key := fmt.Sprintf("blank_%d", i+1)
		inner := match[1]
		parts := strings.Split(inner, ";")
		expected := strings.TrimSpace(parts[0])

		actual := userBlanks[key]
		if strings.TrimSpace(strings.ToLower(actual)) == strings.TrimSpace(strings.ToLower(expected)) {
			correctCount++
		}
	}

	awarded := int(math.Round(float64(correctCount) / float64(len(matches)) * float64(maxPoints)))
	isCorrect := (correctCount == len(matches))

	return domain.BlockValidationResult{
		IsCorrect: isCorrect,
		Feedback:  fmt.Sprintf("Верно %d из %d", correctCount, len(matches)),
	}, awarded
}

func validateInputBlankBlock(props map[string]any, studentAnswer any, maxPoints int) (domain.BlockValidationResult, int) {
	template, _ := props["templateText"].(string)
	matches := dropdownRegex.FindAllStringSubmatch(template, -1)

	userBlanks := extractStringMap(studentAnswer)
	correctCount := 0
	total := len(matches)

	if total > 0 {
		for i, match := range matches {
			key := fmt.Sprintf("blank_%d", i+1)
			expected := strings.TrimSpace(match[1])
			actual := userBlanks[key]
			if strings.TrimSpace(strings.ToLower(actual)) == strings.TrimSpace(strings.ToLower(expected)) {
				correctCount++
			}
		}
	} else if expectedAnswer, ok := props["correctAnswer"].(string); ok && expectedAnswer != "" {
		total = 1
		for _, v := range userBlanks {
			if strings.TrimSpace(strings.ToLower(v)) == strings.TrimSpace(strings.ToLower(expectedAnswer)) {
				correctCount++
				break
			}
		}
	} else {
		return domain.BlockValidationResult{IsCorrect: true}, maxPoints
	}

	awarded := int(math.Round(float64(correctCount) / float64(total) * float64(maxPoints)))
	isCorrect := (correctCount == total)

	return domain.BlockValidationResult{
		IsCorrect: isCorrect,
		Feedback:  fmt.Sprintf("Верно %d из %d", correctCount, total),
	}, awarded
}

func validateSequenceBlock(props map[string]any, studentAnswer any) domain.BlockValidationResult {
	items, ok := props["items"].([]any)
	if !ok || len(items) == 0 {
		return domain.BlockValidationResult{IsCorrect: true}
	}

	expectedOrder := make([]string, 0, len(items))
	for _, it := range items {
		if itMap, ok := it.(map[string]any); ok {
			if t, ok := itMap["text"].(string); ok {
				expectedOrder = append(expectedOrder, strings.TrimSpace(t))
			}
		}
	}

	actualOrder := extractStringSlice(studentAnswer)
	if len(actualOrder) != len(expectedOrder) {
		return domain.BlockValidationResult{
			IsCorrect:     false,
			Feedback:      "Неверный порядок элементов",
			CorrectAnswer: expectedOrder,
		}
	}

	for i := range expectedOrder {
		if strings.TrimSpace(actualOrder[i]) != expectedOrder[i] {
			return domain.BlockValidationResult{
				IsCorrect:     false,
				Feedback:      "Неверный порядок элементов",
				CorrectAnswer: expectedOrder,
			}
		}
	}

	return domain.BlockValidationResult{
		IsCorrect:     true,
		Feedback:      "Верная последовательность",
		CorrectAnswer: expectedOrder,
	}
}

func isBoolCorrect(val any) bool {
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true")
	}
	return false
}

func formatBoolFeedback(ok bool) string {
	if ok {
		return "Верный ответ"
	}
	return "Неверный ответ"
}

func extractIntOption(ans any) *int {
	if ans == nil {
		return nil
	}
	switch v := ans.(type) {
	case float64:
		i := int(v)
		return &i
	case int:
		return &v
	case map[string]any:
		if opt, ok := v["selected_option"]; ok {
			return extractIntOption(opt)
		}
		if opt, ok := v["selectedOption"]; ok {
			return extractIntOption(opt)
		}
	}
	return nil
}

func extractIntSlice(ans any) []int {
	var res []int
	if ans == nil {
		return res
	}

	switch v := ans.(type) {
	case []any:
		for _, it := range v {
			if n := extractIntOption(it); n != nil {
				res = append(res, *n)
			}
		}
	case map[string]any:
		if opts, ok := v["selected_options"]; ok {
			return extractIntSlice(opts)
		}
		if opts, ok := v["selectedOptions"]; ok {
			return extractIntSlice(opts)
		}
	}
	return res
}

func extractStringMap(ans any) map[string]string {
	res := make(map[string]string)
	if ans == nil {
		return res
	}

	switch v := ans.(type) {
	case map[string]string:
		return v
	case map[string]any:
		if pairs, ok := v["pairs"].(map[string]any); ok {
			for pk, pv := range pairs {
				res[pk] = fmt.Sprint(pv)
			}
			return res
		}
		if blanks, ok := v["blanks"].(map[string]any); ok {
			for bk, bv := range blanks {
				res[bk] = fmt.Sprint(bv)
			}
			return res
		}
		for k, val := range v {
			if k != "type" {
				res[k] = fmt.Sprint(val)
			}
		}
	}
	return res
}

func extractStringSlice(ans any) []string {
	var res []string
	if ans == nil {
		return res
	}

	switch v := ans.(type) {
	case []any:
		for _, it := range v {
			res = append(res, fmt.Sprint(it))
		}
	case []string:
		return v
	case map[string]any:
		if ord, ok := v["order"]; ok {
			return extractStringSlice(ord)
		}
	}
	return res
}
