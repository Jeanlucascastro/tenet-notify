package main

import (
	"context"
	"log"
	"tenet-notify/config"
	"tenet-notify/internal/infrastructure"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize context
	ctx := context.Background()

	// Initialize FCM Adapter
	// Note: We're using the path from config, make sure serviceAccountKey.json exists or is pointed to correctly
	fcmAdapter, err := infrastructure.NewFCMAdapter(ctx, cfg.FCMCredentialsPath)
	if err != nil {
		log.Printf("Failed to initialize FCM adapter: %v. Continuing without FCM for testing if credentials missing.", err)
		// In production, we might want to panic or exit, but for development/testing without real creds, we can log.
		// For this task, we assume the user might not have the file yet, so we'll log fatal if we want strictness.
		// Let's implement a 'Mock' or 'NoOp' adapter if file missing? Or just Fail.
		// The prompt says "implementar... usando FCM", implying it should work.
		// I will log fatal to clear that it's required.
		log.Fatalf("Critical error: %v", err)
	} else {
		log.Println("FCM Adapter initialized successfully")
	}

	// Initialize RabbitMQ Consumer
	// We inject the FCM adapter as the NotificationAdapter
	consumer, err := infrastructure.NewRabbitMQConsumer(cfg.RabbitMQURL, fcmAdapter)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer consumer.Close()

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages...")

	// Start consuming (blocking)
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
