package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

type mockNotificationRepo struct {
	items       map[int64]*domain.Notification
	nextID      int64
	unreadCount int64
}

func newMockNotificationRepo() *mockNotificationRepo {
	return &mockNotificationRepo{
		items:  make(map[int64]*domain.Notification),
		nextID: 1,
	}
}

func (m *mockNotificationRepo) CreateNotification(ctx context.Context, n *domain.Notification) (*domain.Notification, error) {
	copyN := *n
	copyN.ID = m.nextID
	m.nextID++
	copyN.CreatedAt = time.Now()
	m.items[copyN.ID] = &copyN
	if !copyN.IsRead {
		m.unreadCount++
	}
	return &copyN, nil
}

func (m *mockNotificationRepo) GetUserNotifications(ctx context.Context, userID int64, limit int) ([]domain.Notification, error) {
	result := make([]domain.Notification, 0)
	for _, item := range m.items {
		if item.UserID == userID {
			result = append(result, *item)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *mockNotificationRepo) CountUnreadNotifications(ctx context.Context, userID int64) (int64, error) {
	var count int64
	for _, item := range m.items {
		if item.UserID == userID && !item.IsRead {
			count++
		}
	}
	return count, nil
}

func (m *mockNotificationRepo) MarkAsRead(ctx context.Context, userID, notificationID int64) error {
	item, ok := m.items[notificationID]
	if !ok || item.UserID != userID {
		return errorsAPP.ErrNotificationNotFound
	}
	item.IsRead = true
	return nil
}

func (m *mockNotificationRepo) MarkAllAsRead(ctx context.Context, userID int64) error {
	for _, item := range m.items {
		if item.UserID == userID {
			item.IsRead = true
		}
	}
	return nil
}

func TestCreateNotification_Success(t *testing.T) {
	repo := newMockNotificationRepo()
	svc := NewNotificationService(repo)

	link := "/lessons/42"
	n, err := svc.CreateNotification(
		context.Background(),
		10,
		"Homework Graded",
		"Your submission was graded: 100 points",
		domain.NotificationTypeHomeworkGraded,
		&link,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if n.ID != 1 {
		t.Errorf("expected ID 1, got %d", n.ID)
	}
	if n.Title != "Homework Graded" {
		t.Errorf("expected title 'Homework Graded', got '%s'", n.Title)
	}
	if n.IsRead {
		t.Errorf("expected IsRead to be false")
	}
	if n.LinkURL == nil || *n.LinkURL != link {
		t.Errorf("expected link %s, got %v", link, n.LinkURL)
	}
}

func TestCreateNotification_ValidationErrors(t *testing.T) {
	repo := newMockNotificationRepo()
	svc := NewNotificationService(repo)

	_, err := svc.CreateNotification(context.Background(), 0, "Title", "Message", domain.NotificationTypeNewLesson, nil)
	if !errors.Is(err, errorsAPP.ErrValidationFailed) {
		t.Errorf("expected ErrValidationFailed for zero userID, got: %v", err)
	}

	_, err = svc.CreateNotification(context.Background(), 1, "", "Message", domain.NotificationTypeNewLesson, nil)
	if !errors.Is(err, errorsAPP.ErrValidationFailed) {
		t.Errorf("expected ErrValidationFailed for empty title, got: %v", err)
	}

	_, err = svc.CreateNotification(context.Background(), 1, "Title", "", domain.NotificationTypeNewLesson, nil)
	if !errors.Is(err, errorsAPP.ErrValidationFailed) {
		t.Errorf("expected ErrValidationFailed for empty message, got: %v", err)
	}
}

func TestGetFeed_Success(t *testing.T) {
	repo := newMockNotificationRepo()
	svc := NewNotificationService(repo)

	_, _ = svc.CreateNotification(context.Background(), 10, "Title 1", "Msg 1", domain.NotificationTypeNewLesson, nil)
	_, _ = svc.CreateNotification(context.Background(), 10, "Title 2", "Msg 2", domain.NotificationTypeHomeworkGraded, nil)
	_, _ = svc.CreateNotification(context.Background(), 99, "Other user", "Msg", domain.NotificationTypeNewLesson, nil)

	feed, err := svc.GetFeed(context.Background(), 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(feed.Notifications) != 2 {
		t.Errorf("expected 2 notifications for user 10, got %d", len(feed.Notifications))
	}
	if feed.UnreadCount != 2 {
		t.Errorf("expected 2 unread, got %d", feed.UnreadCount)
	}
}

func TestMarkAsRead_SuccessAndNotFound(t *testing.T) {
	repo := newMockNotificationRepo()
	svc := NewNotificationService(repo)

	n, _ := svc.CreateNotification(context.Background(), 10, "Title", "Msg", domain.NotificationTypeNewLesson, nil)

	err := svc.MarkAsRead(context.Background(), 10, n.ID)
	if err != nil {
		t.Fatalf("unexpected error marking as read: %v", err)
	}

	feed, _ := svc.GetFeed(context.Background(), 10, 10)
	if feed.UnreadCount != 0 {
		t.Errorf("expected 0 unread, got %d", feed.UnreadCount)
	}

	// Not found
	err = svc.MarkAsRead(context.Background(), 10, 999)
	if !errors.Is(err, errorsAPP.ErrNotificationNotFound) {
		t.Errorf("expected ErrNotificationNotFound, got: %v", err)
	}

	// Another user trying to read user 10's notification
	err = svc.MarkAsRead(context.Background(), 88, n.ID)
	if !errors.Is(err, errorsAPP.ErrNotificationNotFound) {
		t.Errorf("expected ErrNotificationNotFound for wrong user, got: %v", err)
	}
}

func TestMarkAllAsRead_Success(t *testing.T) {
	repo := newMockNotificationRepo()
	svc := NewNotificationService(repo)

	_, _ = svc.CreateNotification(context.Background(), 10, "Title 1", "Msg 1", domain.NotificationTypeNewLesson, nil)
	_, _ = svc.CreateNotification(context.Background(), 10, "Title 2", "Msg 2", domain.NotificationTypeNewSubmission, nil)

	feedBefore, _ := svc.GetFeed(context.Background(), 10, 10)
	if feedBefore.UnreadCount != 2 {
		t.Fatalf("expected 2 unread before, got %d", feedBefore.UnreadCount)
	}

	err := svc.MarkAllAsRead(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error marking all as read: %v", err)
	}

	feedAfter, _ := svc.GetFeed(context.Background(), 10, 10)
	if feedAfter.UnreadCount != 0 {
		t.Errorf("expected 0 unread after, got %d", feedAfter.UnreadCount)
	}
}
