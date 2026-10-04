package domain

import "time"

type NotificationType string

const (
	NotificationTypeHomeworkGraded NotificationType = "homework_graded"
	NotificationTypeNewSubmission  NotificationType = "new_submission"
	NotificationTypeNewLesson      NotificationType = "new_lesson"
)

type Notification struct {
	ID        int64            `json:"id"`
	UserID    int64            `json:"user_id"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Type      NotificationType `json:"type"`
	LinkURL   *string          `json:"link_url,omitempty"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}

type NotificationFeed struct {
	Notifications []Notification `json:"notifications"`
	UnreadCount   int64          `json:"unread_count"`
}
