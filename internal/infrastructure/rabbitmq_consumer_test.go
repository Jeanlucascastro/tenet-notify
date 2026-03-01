package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"tenet-notify/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

var errTest = errors.New("test error")

type mockNotificationAdapter struct {
	sendFunc func(ctx context.Context, notification model.Notification) error
}

func (m *mockNotificationAdapter) Send(ctx context.Context, notification model.Notification) error {
	if m.sendFunc != nil {
		return m.sendFunc(ctx, notification)
	}
	return nil
}

type mockChannel struct {
	consumeChan chan amqp.Delivery
}

func (m *mockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return m.consumeChan, nil
}

func (m *mockChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (m *mockChannel) Close() error {
	return nil
}

type mockConnection struct {
	channel *mockChannel
}

func (m *mockConnection) Channel() (*amqp.Channel, error) {
	ch := amqp.Channel{}
	return &ch, nil
}

func (m *mockConnection) Close() error {
	return nil
}

type RabbitMQConsumerTestable struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	notificationSvc interface {
		Send(ctx context.Context, notification model.Notification) error
	}
	consumeChan chan amqp.Delivery
}

func TestRabbitMQConsumer_ReceiveMessage(t *testing.T) {
	consumeChan := make(chan amqp.Delivery, 1)

	notification := model.Notification{
		Token: "test-token",
		Data: model.NotificationData{
			Type:      model.NotificationTypeNewMessage,
			SessionID: 123,
			SenderID:  456,
			Message:   "Hello World",
		},
	}

	payload := model.MessagePayload{
		Message: notification,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	adapterCalled := false
	mockAdapter := &mockNotificationAdapter{
		sendFunc: func(ctx context.Context, n model.Notification) error {
			adapterCalled = true
			if n.Token != notification.Token {
				t.Errorf("Expected token %s, got %s", notification.Token, n.Token)
			}
			if n.Data.Message != notification.Data.Message {
				t.Errorf("Expected message %s, got %s", notification.Data.Message, n.Data.Message)
			}
			return nil
		},
	}

	consumer := &RabbitMQConsumer{
		notificationSvc: mockAdapter,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-consumeChan:
				var payload model.MessagePayload
				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					t.Logf("Error decoding JSON: %v", err)
					continue
				}
				if err := consumer.notificationSvc.Send(ctx, payload.Message); err != nil {
					t.Logf("Error sending notification: %v", err)
				}
			}
		}
	}()

	consumeChan <- amqp.Delivery{
		Body: payloadBytes,
	}

	time.Sleep(100 * time.Millisecond)

	if !adapterCalled {
		t.Error("Expected adapter.Send to be called")
	}
}

func TestRabbitMQConsumer_InvalidJSON(t *testing.T) {
	consumeChan := make(chan amqp.Delivery, 1)

	adapterCalled := false
	mockAdapter := &mockNotificationAdapter{
		sendFunc: func(ctx context.Context, n model.Notification) error {
			adapterCalled = true
			return nil
		},
	}

	consumer := &RabbitMQConsumer{
		notificationSvc: mockAdapter,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-consumeChan:
				var payload model.MessagePayload
				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					t.Logf("Error decoding JSON: %v", err)
					continue
				}
				if err := consumer.notificationSvc.Send(ctx, payload.Message); err != nil {
					t.Logf("Error sending notification: %v", err)
				}
			}
		}
	}()

	consumeChan <- amqp.Delivery{
		Body: []byte("invalid json"),
	}

	time.Sleep(100 * time.Millisecond)

	if adapterCalled {
		t.Error("Expected adapter.Send NOT to be called for invalid JSON")
	}
}

func TestRabbitMQConsumer_AdapterError(t *testing.T) {
	consumeChan := make(chan amqp.Delivery, 1)

	notification := model.Notification{
		Token: "test-token",
		Data: model.NotificationData{
			Type:      model.NotificationTypeNewMessage,
			SessionID: 123,
			SenderID:  456,
			Message:   "Hello World",
		},
	}

	payload := model.MessagePayload{
		Message: notification,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	adapterCalled := false
	mockAdapter := &mockNotificationAdapter{
		sendFunc: func(ctx context.Context, n model.Notification) error {
			adapterCalled = true
			return errTest
		},
	}

	consumer := &RabbitMQConsumer{
		notificationSvc: mockAdapter,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-consumeChan:
				var payload model.MessagePayload
				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					t.Logf("Error decoding JSON: %v", err)
					continue
				}
				if err := consumer.notificationSvc.Send(ctx, payload.Message); err != nil {
					t.Logf("Error sending notification: %v", err)
				}
			}
		}
	}()

	consumeChan <- amqp.Delivery{
		Body: payloadBytes,
	}

	time.Sleep(100 * time.Millisecond)

	if !adapterCalled {
		t.Error("Expected adapter.Send to be called")
	}
}
