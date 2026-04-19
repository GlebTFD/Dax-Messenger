package dto

// type Message struct {
// 	ID        string
// 	Type      string
// 	Timestamp int64
// 	Text      string
// 	ReplyTo   string
// }

type MessageJSON struct {
	ID        string             `json:"id"`
	Type      string             `json:"type"`
	Timestamp int64              `json:"timestamp"`
	Payload   TextMessagePayload `json:"payload"`
}

type TextMessagePayload struct {
	Text    string `json:"text"`
	ReplyTo string `json:"replyTo"`
}

type PubSubMessage struct {
	Channel string
	Payload string
}

// RedisMessage is a unified structure for all pub/sub message types.
// Type can be "message", "message_deleted", or "message_updated".
type RedisMessage struct {
	Channel string      `json:"channel"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type DeleteMessageResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type DeletedNotificationPayload struct {
	ID string `json:"id"`
}

type DeleteNotification struct {
	Type    string                     `json:"type"`
	Payload DeletedNotificationPayload `json:"payload"`
}

type UpdateMessageRequest struct {
	Text string `json:"text"`
}

type UpdateMessageResponse struct {
	ID      string `json:"id"`
	Updated bool   `json:"updated"`
	Text    string `json:"text"`
}

type UpdatedNotificationPayload struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type UpdateNotification struct {
	Type    string                     `json:"type"`
	Payload UpdatedNotificationPayload `json:"payload"`
}
