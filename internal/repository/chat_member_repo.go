package repository

import (
	"go-chat/internal/domain"
	"go-chat/pkg/db"
)

// 取得聊天室成員的 userID 列表
func GetUserIDsByChatID(chatID uint) ([]uint, error) {
	var userIDs []uint

	err := db.DB.
		Table("chat_members").
		Where("chat_id = ?", chatID).
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}

// 檢查 userID 是否為 chatID 的成員
func IsChatMember(userID, chatID uint) (bool, error) {
	var count int64
	err := db.DB.
		Table("chat_members").
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Count(&count).Error
	return count > 0, err
}

// 取得聊天室成員的 user_id 和 username 列表
func GetGroupMembers(chatID uint) ([]domain.MemberInfo, error) {
	var members []domain.MemberInfo
	err := db.DB.
		Table("chat_members cm").
		Select("u.id as user_id, u.username").
		Joins("JOIN users u ON u.id = cm.user_id").
		Where("cm.chat_id = ?", chatID).
		Scan(&members).Error
	return members, err
}

// 將使用者加入聊天室
func AddMemberToChat(chatID, userID uint) error {
	return db.DB.Create(&domain.ChatMember{ChatID: chatID, UserID: userID}).Error
}
