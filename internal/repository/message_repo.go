package repository

import (
	"go-chat/internal/domain"
	"go-chat/pkg/db"
)

func SaveMessage(msg *domain.Message) error {
	return db.DB.Create(msg).Error
}
func GetUnreadMessages(userID string) ([]domain.Message, error) {
	var msgs []domain.Message
	err := db.DB.Where("target_id = ? AND read = ?", userID, false).Find(&msgs).Error
	return msgs, err
}
func MarkAsRead(userID string) {
	db.DB.Model(&domain.Message{}).
		Where("target_id = ? AND read = ?", userID, false).
		Update("read", true)
}