package repository

import "go-chat/pkg/db"

// UpsertMessageRead records or updates the last read message for a user in a chat.
// Because the message_reads table has a FK on id referencing messages(id),
// we store lastReadMessageID in the id column on first insert to satisfy that FK.
// On conflict we only update last_read_message_id and read_at (not the PK).
func UpsertMessageRead(userID, chatID, lastReadMessageID uint) error {
	return db.DB.Exec(`
		INSERT INTO message_reads (id, user_id, chat_id, last_read_message_id, read_at)
		VALUES (?, ?, ?, ?, NOW())
		ON CONFLICT (user_id, chat_id)
		DO UPDATE SET last_read_message_id = ?, read_at = NOW()
	`, lastReadMessageID, userID, chatID, lastReadMessageID,
		lastReadMessageID).Error
}
