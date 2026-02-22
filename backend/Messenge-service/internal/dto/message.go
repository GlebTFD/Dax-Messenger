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
