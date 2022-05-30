package quote

import (
	"time"
)

type Quote struct {
	tableName     struct{}  `pg:"quote_quote"`
	ID            int64     `json:"id"`
	Text          string    `json:"text"`
	UserID        int64     `json:"userId"`
	CreatedAt     time.Time `json:"created_at"`
	MessageSentAt time.Time `json:"message_sent_at"`
}
