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
