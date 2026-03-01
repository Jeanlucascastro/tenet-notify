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
	SessionID int64            `json:"sessionId"`
	SenderID  int64            `json:"senderId"`
	Message   string           `json:"message"`
}

type MessagePayload struct {
	Message Notification `json:"message"`
}
