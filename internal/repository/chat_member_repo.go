package repository

import (
	"go-chat/internal/domain"
	"go-chat/pkg/db"
)

func GetUserIDsByChatID(chatID uint) ([]uint, error) {
	var userIDs []uint

	err := db.DB.
		Table("chat_members").
		Where("chat_id = ?", chatID).
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}

// IsChatMember returns true if userID is a member of chatID.
func IsChatMember(userID, chatID uint) (bool, error) {
	var count int64
	err := db.DB.
		Table("chat_members").
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetGroupMembers returns id and username of every member in the given chat.
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

// AddMemberToChat inserts a new member row for the given chat and user.
func AddMemberToChat(chatID, userID uint) error {
	return db.DB.Create(&domain.ChatMember{ChatID: chatID, UserID: userID}).Error
}
