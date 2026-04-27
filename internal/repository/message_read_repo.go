package repository

import "go-chat/pkg/db"

// UpsertMessageRead records or updates the last read message for a user in a chat.
// The table uses (user_id, chat_id) as its composite primary key — there is no
// surrogate id column.
func UpsertMessageRead(userID, chatID, lastReadMessageID uint) error {
	return db.DB.Exec(`
		INSERT INTO message_reads (user_id, chat_id, last_read_message_id, read_at)
		VALUES (?, ?, ?, NOW())
		ON CONFLICT (user_id, chat_id)
		DO UPDATE SET last_read_message_id = ?, read_at = NOW()
	`, userID, chatID, lastReadMessageID, lastReadMessageID).Error
}
