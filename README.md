# Tenet Notify Microservice

Microservice for handling notifications using RabbitMQ and Firebase Cloud Messaging (FCM).

## Architecture

- **Core Domain**: Defines `Notification` model.
- **Adapters**:
  - `FCMAdapter`: Sends notifications via Firebase.
  - `RabbitMQConsumer`: Listens for messages on `notifications` queue.
- **Config**: Loads environment variables.

## Prerequisites

- Go 1.26+
- RabbitMQ instance
- Firebase Project with `serviceAccountKey.json`

## Configuration

Set the following environment variables (or use defaults):

- `RABBITMQ_URL`: AMQP connection string (default: `amqp://guest:guest@localhost:5672/`)
- `FCM_CREDENTIALS_PATH`: Path to Firebase credentials file (default: `serviceAccountKey.json`)

## Running

1. Place your `serviceAccountKey.json` in the root directory.
2. Run the application:
   ```bash
   go run cmd/main.go
   ```

## Message Format

Send JSON messages to the `notifications` queue:

```json
{
  "message": {
    "token": "FCM_TOKEN_DO_USUARIO",
    "data": {
      "type": "NEW_MESSAGE",
      "sessionId": "123",
      "senderId": "45",
      "message": "Oi, tudo bem?"
    }
  }
}
```
