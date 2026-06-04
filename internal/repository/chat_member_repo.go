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

// AddMemberToChat adds a user to a chat if not already a member.
func AddMemberToChat(chatID, userID uint) error {
	isMember, err := IsChatMember(userID, chatID)
	if err != nil {
		return err
	}
	if isMember {
		return nil
	}
	return db.DB.Create(&domain.ChatMember{ChatID: chatID, UserID: userID}).Error
}
