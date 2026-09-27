package analytics

import (
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

func (h *AnalyticsHandler) GetStudentDrilldown(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.analytics.GetStudentDrilldown"

	log := logger.GetLogger(r.Context(), op)

	courseIDStr := chi.URLParam(r, "courseid")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	studentIDStr := chi.URLParam(r, "userid")
	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	teacherID := h.authMiddleware.GetUserID(r.Context())
	if teacherID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	report, err := h.analyticsService.GetStudentDrilldown(r.Context(), teacherID, courseID, studentID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	lessonLogs := make([]dto.StudentLessonLogDTO, 0, len(report.LessonLogs))
	for _, l := range report.LessonLogs {
		lessonLogs = append(lessonLogs, dto.StudentLessonLogDTO{
			LessonID:       l.LessonID,
			LessonTitle:    l.LessonTitle,
			LessonType:     l.LessonType,
			Status:         l.Status,
			Score:          l.Score,
			CompletedAt:    l.CompletedAt,
			LastAccessedAt: l.LastAccessedAt,
		})
	}

	testAttempts := make([]dto.StudentTestAttemptDTO, 0, len(report.TestAttempts))
	for _, att := range report.TestAttempts {
		answers := make([]dto.StudentAttemptAnswerDetailDTO, 0, len(att.Answers))
		for _, ans := range att.Answers {
			answers = append(answers, dto.StudentAttemptAnswerDetailDTO{
				QuestionID:    ans.QuestionID,
				QuestionText:  ans.QuestionText,
				ChosenAnswer:  ans.ChosenAnswer,
				IsCorrect:     ans.IsCorrect,
				Points:        ans.Points,
				MaxPoints:     ans.MaxPoints,
				Feedback:      ans.Feedback,
				AttachmentURL: ans.AttachmentURL,
				Explanation:   ans.Explanation,
			})
		}

		testAttempts = append(testAttempts, dto.StudentTestAttemptDTO{
			AttemptID:     att.AttemptID,
			QuizID:        att.QuizID,
			QuizTitle:     att.QuizTitle,
			AttemptNumber: att.AttemptNumber,
			Score:         att.Score,
			Passed:        att.Passed,
			StartedAt:     att.StartedAt,
			CompletedAt:   att.CompletedAt,
			Answers:       answers,
		})
	}

	studentName := report.Student.Username
	if report.Student.FirstName != nil && report.Student.LastName != nil {
		studentName = *report.Student.FirstName + " " + *report.Student.LastName
	}

	resp := dto.StudentDrilldownResponse{
		StudentID:       report.Student.ID,
		StudentName:     studentName,
		StudentEmail:    report.Student.Email,
		StudentUsername: report.Student.Username,
		AvatarURL:       report.Student.AvatarURL,
		Bio:             report.Student.Bio,
		CourseID:        report.CourseID,
		OverallProgress: report.OverallProgress,
		AvgScore:        report.AvgScore,
		LessonLogs:      lessonLogs,
		TestAttempts:    testAttempts,
	}

	response.OK(w, r, resp)
}
