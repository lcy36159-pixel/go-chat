package repository

import "go-chat/pkg/db"

func GetUserIDsByChatID(chatID uint) ([]uint, error) {
	var userIDs []uint

	err := db.DB.
		Table("chat_members").
		Where("chat_id = ?", chatID).
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}
