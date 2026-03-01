package infrastructure

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"tenet-notify/internal/model"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FCMAdapter struct {
	client *messaging.Client
}

func NewFCMAdapter(ctx context.Context, credentialsPath string) (*FCMAdapter, error) {
	opt := option.WithAuthCredentialsFile(option.ServiceAccount, credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %v", err)
	}

	return &FCMAdapter{client: client}, nil
}

func (a *FCMAdapter) Send(ctx context.Context, notification model.Notification) error {
	message := &messaging.Message{
		Token: notification.Token,
		Data: map[string]string{
			"type":      string(notification.Data.Type),
			"sessionId": strconv.FormatInt(notification.Data.SessionID, 10),
			"senderId":  strconv.FormatInt(notification.Data.SenderID, 10),
			"message":   notification.Data.Message,
		},
	}

	response, err := a.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	log.Printf("Successfully sent message: %s", response)
	return nil
}
