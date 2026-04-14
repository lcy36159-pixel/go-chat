package service

import (
	"errors"
	"go-chat/internal/domain"
	"go-chat/internal/repository"
)

// 處理訊息（核心入口）
func HandleMessage(msg *domain.Message) error {
	// 基本驗證
	if msg.ChatID == 0 {
		return errors.New("chat_id is required")
	}
	if msg.SenderID == nil {
		return errors.New("sender_id is required")
	}
	if msg.Content == "" {
		return errors.New("content is empty")
	}

	// 預設 type（避免漏傳）
	if msg.Type == "" {
		msg.Type = "text"
	}

	// 可以在這裡加更多邏輯（之後擴充）
	// ex: 長度限制、過濾字詞、圖片處理...

	return repository.SaveMessage(msg)
}

// 建立群組聊天室
func CreateGroupChat(creatorID uint, name string, userIDs []uint) (uint, error) {
	// 1️⃣ 建立 chat
	chat := domain.Chat{
		Name:      name,
		Type:      "group",
		CreatedBy: creatorID,
	}

	if err := repository.CreateChat(&chat); err != nil {
		return 0, err
	}

	// 2️⃣ 加入成員（包含自己）
	memberMap := make(map[uint]bool)
	memberMap[creatorID] = true

	for _, id := range userIDs {
		memberMap[id] = true
	}

	for id := range memberMap {
		repository.AddChatMember(&domain.ChatMember{
			ChatID: chat.ID,
			UserID: id,
		})
	}

	return chat.ID, nil
}
