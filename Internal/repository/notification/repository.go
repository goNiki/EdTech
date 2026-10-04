package notification

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type notificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) repository.NotificationRepository {
	return &notificationRepository{pool: pool}
}

func (r *notificationRepository) CreateNotification(ctx context.Context, n *domain.Notification) (*domain.Notification, error) {
	const op = "repository.notification.CreateNotification"
	q := txmanager.GetQueryExecutor(ctx, r.pool)

	query := `
		INSERT INTO notifications (user_id, title, message, type, link_url, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, COALESCE(NULLIF($7, '0001-01-01 00:00:00+00'::timestamptz), NOW()))
		RETURNING id, created_at
	`

	res := *n
	err := q.QueryRow(ctx, query,
		res.UserID,
		res.Title,
		res.Message,
		string(res.Type),
		res.LinkURL,
		res.IsRead,
		res.CreatedAt,
	).Scan(&res.ID, &res.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return &res, nil
}

func (r *notificationRepository) GetUserNotifications(ctx context.Context, userID int64, limit int) ([]domain.Notification, error) {
	const op = "repository.notification.GetUserNotifications"
	q := txmanager.GetQueryExecutor(ctx, r.pool)

	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT id, user_id, title, message, type, link_url, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := q.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	notifications := make([]domain.Notification, 0)
	for rows.Next() {
		var n domain.Notification
		var typeStr string
		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Title,
			&n.Message,
			&typeStr,
			&n.LinkURL,
			&n.IsRead,
			&n.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan error: %w", op, err)
		}
		n.Type = domain.NotificationType(typeStr)
		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return notifications, nil
}

func (r *notificationRepository) CountUnreadNotifications(ctx context.Context, userID int64) (int64, error) {
	const op = "repository.notification.CountUnreadNotifications"
	q := txmanager.GetQueryExecutor(ctx, r.pool)

	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = FALSE
	`

	var count int64
	if err := q.QueryRow(ctx, query, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return count, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, userID, notificationID int64) error {
	const op = "repository.notification.MarkAsRead"
	q := txmanager.GetQueryExecutor(ctx, r.pool)

	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE id = $1 AND user_id = $2
	`

	tag, err := q.Exec(ctx, query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotificationNotFound)
	}

	return nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	const op = "repository.notification.MarkAllAsRead"
	q := txmanager.GetQueryExecutor(ctx, r.pool)

	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE user_id = $1 AND is_read = FALSE
	`

	_, err := q.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return nil
}
