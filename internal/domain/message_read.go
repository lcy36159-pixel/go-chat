package domain

import "time"

// MessageRead tracks the last read message for a user in a chat.
// It maps to the message_reads table (also referred to as chat_reads).
type MessageRead struct {
	ID                uint      `gorm:"primaryKey"`
	UserID            uint      `gorm:"not null;uniqueIndex:idx_user_chat"`
	ChatID            uint      `gorm:"not null;uniqueIndex:idx_user_chat"`
	LastReadMessageID uint      `gorm:"not null"`
	ReadAt            time.Time `gorm:"autoUpdateTime"`
}
