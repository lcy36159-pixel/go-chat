package repository

import (
	"go-chat/internal/domain"
	"go-chat/pkg/db"
)

func SaveMessage(msg *domain.Message) error {
	return db.DB.Create(msg).Error
}

func GetMessagesByChatID(chatID uint, lastID uint, limit int) ([]domain.Message, error) {
	var msgs []domain.Message

	query := db.DB.
		Where("chat_id = ?", chatID).
		Order("id DESC").
		Limit(limit)

	if lastID > 0 {
		query = query.Where("id < ?", lastID)
	}

	err := query.Find(&msgs).Error
	return msgs, err
}
