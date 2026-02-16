package service

import (
	"context"
	"tenet-notify/internal/model"
)

// NotificationAdapter defines the interface for sending notifications via different providers (FCM, APNS, etc.)
type NotificationAdapter interface {
	Send(ctx context.Context, notification model.Notification) error
}

// NotificationService defines the interface for the business logic of processing notifications
type NotificationService interface {
	ProcessNotification(ctx context.Context, payload model.MessagePayload) error
}
