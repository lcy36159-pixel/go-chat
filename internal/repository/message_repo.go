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
func GetMessages(user1, user2 string, lastID uint) ([]domain.Message, error) {
	var msgs []domain.Message

	query := db.DB.
		Where(
			"(sender_id = ? AND target_id = ?) OR (sender_id = ? AND target_id = ?)",
			user1, user2, user2, user1,
		).
		Order("id DESC").
		Limit(20)

	if lastID > 0 {
		query = query.Where("id < ?", lastID)
	}

	err := query.Find(&msgs).Error
	return msgs, err
}
