package model

type NotificationType string

const (
	NotificationTypeNewMessage NotificationType = "NEW_MESSAGE"
)

type Notification struct {
	Token string           `json:"token"`
	Data  NotificationData `json:"data"`
}

type NotificationData struct {
	Type      NotificationType `json:"type"`
	SessionID string           `json:"sessionId"`
	SenderID  string           `json:"senderId"`
	Message   string           `json:"message"`
}

type MessagePayload struct {
	Message Notification `json:"message"`
}
