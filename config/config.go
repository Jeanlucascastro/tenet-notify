package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL        string
	FCMCredentialsPath string
	FirebaseProjectID  string
}

func Load() *Config {

	// Carrega .env apenas se existir (não quebra em produção)
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	vhost := os.Getenv("RABBITMQ_VHOST")

	print("user ----" + user)

	rabbitURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s%s",
		user,
		pass,
		host,
		port,
		vhost,
	)

	fcmCredentialsPath := os.Getenv("FCM_CREDENTIALS_PATH")
	if fcmCredentialsPath == "" {
		fcmCredentialsPath = "tenet-9739c-firebase-adminsdk-fbsvc-c87629ef6b.json"
	}

	firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID")
	if firebaseProjectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID is required")
	}

	return &Config{
		RabbitMQURL:        rabbitURL,
		FCMCredentialsPath: fcmCredentialsPath,
		FirebaseProjectID:  firebaseProjectID,
	}
}
