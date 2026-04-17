package service

import (
	"errors"
	"go-chat/internal/domain"
	"go-chat/internal/repository"
	"go-chat/pkg/db"
)

// 建立群組聊天室（完整安全版）
func CreateGroupChat(creatorID uint, name string, userIDs []uint) (uint, error) {
	// 基本驗證
	if name == "" {
		return 0, errors.New("name is required")
	}

	// 使用 transaction（關鍵）
	tx := db.DB.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// ❗任何錯誤都 rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1️⃣ 建立 chat
	chat := domain.Chat{
		Name:      name,
		Type:      "group",
		CreatedBy: creatorID,
	}

	if err := tx.Create(&chat).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 2️⃣ 整理 member（去重 + 一定包含自己）
	memberMap := make(map[uint]bool)
	memberMap[creatorID] = true

	for _, id := range userIDs {
		memberMap[id] = true
	}

	// 3️⃣ 寫入 chat_members
	for id := range memberMap {
		member := domain.ChatMember{
			ChatID: chat.ID,
			UserID: id,
		}

		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// 4️⃣ commit
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return chat.ID, nil
}
func GetUserChats(userID uint) ([]domain.ChatDTO, error) {
	return repository.GetUserChats(userID)
}
