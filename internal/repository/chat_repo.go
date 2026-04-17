package repository

import (
	"go-chat/internal/domain"
	"go-chat/pkg/db"
)

// 建立聊天室
func CreateChat(chat *domain.Chat) error {
	return db.DB.Create(chat).Error
}

// 加入成員
func AddChatMember(member *domain.ChatMember) error {
	return db.DB.Create(member).Error
}
func GetUserChats(userID uint) ([]domain.ChatDTO, error) {
	var chats []domain.ChatDTO
	err := db.DB.Table("chat_members").
		Select("chats.id as chat_id, chats.name, messages.content as last_message, messages.created_at as updated_at").
		Joins("left join chats on chats.id = chat_members.chat_id").
		Joins("left join messages on messages.chat_id = chats.id").
		Where("chat_members.user_id = ?", userID).
		Group("chats.id").
		Order("messages.created_at desc").
		Scan(&chats).Error
	return chats, err
}
