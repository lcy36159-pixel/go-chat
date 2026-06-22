package repository

import (
	"fmt"
	"go-chat/internal/domain"
	"go-chat/pkg/db"
)

// 建立聊天室
func CreateChat(chat *domain.Chat) error {
	return db.DB.Create(chat).Error
}

// 取得使用者的聊天室列表
func GetUserChats(userID uint) ([]domain.ChatDTO, error) {
	var chats []domain.ChatDTO

	err := db.DB.
		Table("chats").
		Select(`
			chats.id as chat_id,
			CASE 
				WHEN chats.type = 'private' THEN u.username
				ELSE chats.name
			END as name,
			m.content as last_message,
			m.created_at as updated_at,
			(
				SELECT COUNT(*)
				FROM messages unread
				WHERE unread.chat_id = chats.id
				  AND unread.id > COALESCE(mr.last_read_message_id, 0)
				  AND unread.sender_id != ?
				  AND unread.deleted_at IS NULL
			) as unread_count
		`, userID).
		Joins("JOIN chat_members cm ON cm.chat_id = chats.id").
		Joins(`
			LEFT JOIN chat_members cm2 
			ON cm2.chat_id = chats.id AND cm2.user_id != ? AND chats.type = 'private'
		`, userID).
		Joins(`
			LEFT JOIN users u ON u.id = cm2.user_id
		`).
		Joins(`
			LEFT JOIN messages m ON m.id = (
				SELECT id FROM messages
				WHERE messages.chat_id = chats.id
				ORDER BY created_at DESC
				LIMIT 1
			)
		`).
		Joins(`
			LEFT JOIN message_reads mr ON mr.chat_id = chats.id AND mr.user_id = ?
		`, userID).
		Where("cm.user_id = ?", userID).
		Order("m.created_at DESC NULLS LAST").
		Scan(&chats).Error

	return chats, err
}
func FindPrivateChat(user1, user2 uint) (uint, error) {
	key := GeneratePrivateKey(user1, user2)

	var chatID uint
	err := db.DB.
		Table("chats").
		Select("id").
		Where("search_key = ?", key).
		Limit(1).
		Scan(&chatID).Error

	return chatID, err
}

// GetChatByID returns the chat with the given ID.
func GetChatByID(chatID uint) (*domain.Chat, error) {
	var chat domain.Chat
	err := db.DB.First(&chat, chatID).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func GeneratePrivateKey(user1, user2 uint) string {
	if user1 < user2 {
		return fmt.Sprintf("%d_%d", user1, user2)
	}
	return fmt.Sprintf("%d_%d", user2, user1)
}
