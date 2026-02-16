package infrastructure

import (
	"encoding/json"
	"log"
	"tenet-notify/internal/model"
	"tenet-notify/internal/service"

	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	notificationSvc service.NotificationAdapter // In a real scenario, this might inject a higher-level service, but adapter fits the wire-up for now or we create a concrete service.
}

// For this task, the requirements say "vai chegar uma mensagem pelo rabbit".
// We need to consume it and use the adapter to send it.
// The structure provided in the prompt is:
// {
//   "message": { ... }
// }
// This matches model.MessagePayload.

func NewRabbitMQConsumer(url string, adapter service.NotificationAdapter) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare queue to ensure it exists
	_, err = ch.QueueDeclare(
		"notifications", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:            conn,
		channel:         ch,
		notificationSvc: adapter,
	}, nil
}

func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		"notifications", // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var payload model.MessagePayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			// Use the adapter to send the notification
			// In a more complex app, this would call a domain service method like UseCase.SendNotification
			// which would then choose the adapter. For now, direct adapter usage as per requirements "adapter para facilitar outros tipos...".
			if err := c.notificationSvc.Send(ctx, payload.Message); err != nil {
				log.Printf("Error sending notification: %s", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func (c *RabbitMQConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
