package domain

import "time"

// MessageRead tracks the last read message for a user in a chat.
// PRIMARY KEY is the composite (chat_id, user_id) — no surrogate id column.
type MessageRead struct {
	UserID            uint      `gorm:"primaryKey;not null"`
	ChatID            uint      `gorm:"primaryKey;not null"`
	LastReadMessageID uint      `gorm:"not null"`
	ReadAt            time.Time `gorm:"autoUpdateTime"`
}
