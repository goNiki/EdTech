package notification

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/repository"
	services "edtech/internal/service"
	errorsAPP "edtech/pkg/errors"
)

type service struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) services.NotificationServices {
	return &service{
		notificationRepo: notificationRepo,
	}
}

func (s *service) CreateNotification(
	ctx context.Context,
	userID int64,
	title, message string,
	nType domain.NotificationType,
	linkURL *string,
) (*domain.Notification, error) {
	const op = "service.notification.CreateNotification"

	if userID <= 0 {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrValidationFailed)
	}
	if title == "" || message == "" {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrValidationFailed)
	}

	n := &domain.Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    nType,
		LinkURL: linkURL,
		IsRead:  false,
	}

	created, err := s.notificationRepo.CreateNotification(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return created, nil
}

func (s *service) GetFeed(ctx context.Context, userID int64, limit int) (*domain.NotificationFeed, error) {
	const op = "service.notification.GetFeed"

	if userID <= 0 {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrValidationFailed)
	}

	if limit <= 0 {
		limit = 20
	}

	notifications, err := s.notificationRepo.GetUserNotifications(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	unreadCount, err := s.notificationRepo.CountUnreadNotifications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.NotificationFeed{
		Notifications: notifications,
		UnreadCount:   unreadCount,
	}, nil
}

func (s *service) MarkAsRead(ctx context.Context, userID, notificationID int64) error {
	const op = "service.notification.MarkAsRead"

	if userID <= 0 || notificationID <= 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrValidationFailed)
	}

	if err := s.notificationRepo.MarkAsRead(ctx, userID, notificationID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) MarkAllAsRead(ctx context.Context, userID int64) error {
	const op = "service.notification.MarkAllAsRead"

	if userID <= 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrValidationFailed)
	}

	if err := s.notificationRepo.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
