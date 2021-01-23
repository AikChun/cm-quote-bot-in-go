package telegrambot

type SendMessagePayload struct {
	ChatId           int64  `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageID string `json:"reply_to_message_id", omitempty`
}
