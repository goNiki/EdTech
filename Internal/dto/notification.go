package dto

import (
	"time"

	"edtech/internal/domain"
)

type NotificationResponse struct {
	ID        int64                   `json:"id"`
	UserID    int64                   `json:"user_id"`
	Title     string                  `json:"title"`
	Message   string                  `json:"message"`
	Type      domain.NotificationType `json:"type"`
	LinkURL   *string                 `json:"link_url,omitempty"`
	IsRead    bool                    `json:"is_read"`
	CreatedAt string                  `json:"created_at"`
}

type NotificationFeedResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                  `json:"unread_count"`
}

func ToNotificationResponse(n domain.Notification) NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Message:   n.Message,
		Type:      n.Type,
		LinkURL:   n.LinkURL,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func ToNotificationFeedResponse(feed *domain.NotificationFeed) NotificationFeedResponse {
	if feed == nil {
		return NotificationFeedResponse{
			Notifications: []NotificationResponse{},
			UnreadCount:   0,
		}
	}

	items := make([]NotificationResponse, 0, len(feed.Notifications))
	for _, n := range feed.Notifications {
		items = append(items, ToNotificationResponse(n))
	}

	return NotificationFeedResponse{
		Notifications: items,
		UnreadCount:   feed.UnreadCount,
	}
}
