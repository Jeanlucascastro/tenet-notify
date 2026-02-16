package config

import (
	"fmt"
	"os"
)

type Config struct {
	RabbitMQURL        string
	FCMCredentialsPath string
}

func Load() *Config {
	user := getEnv("RABBITMQ_USER", "guest")
	pass := getEnv("RABBITMQ_PASS", "guest")
	host := getEnv("RABBITMQ_HOST", "localhost")
	port := getEnv("RABBITMQ_PORT", "5672")
	vhost := getEnv("RABBITMQ_VHOST", "/")

	rabbitURL := getEnv(
		"RABBITMQ_URL",
		fmt.Sprintf("amqp://%s:%s@%s:%s%s", user, pass, host, port, vhost),
	)

	return &Config{
		RabbitMQURL:        rabbitURL,
		FCMCredentialsPath: getEnv("FCM_CREDENTIALS_PATH", "serviceAccountKey.json"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
