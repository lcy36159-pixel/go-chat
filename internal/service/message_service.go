package service

import (
	"encoding/json"
	"go-chat/internal/domain"
	"go-chat/internal/repository"
	"go-chat/internal/ws"
)

// 向聊天室發送訊息
func SendMessage(userID, chatID uint, content string) error {
	// 基本驗證
	if userID == 0 {
		return ErrUserIDRequired
	}
	if chatID == 0 {
		return ErrChatIDRequired
	}
	if content == "" {
		return ErrContentEmpty
	}
	// 檢查是否為聊天室成員
	isMember, err := repository.IsChatMember(userID, chatID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrNotChatMember
	}

	message := domain.Message{
		ChatID:   chatID,
		SenderID: &userID,
		Type:     "text",
		Content:  content,
	}
	// 處理訊息 (存入資料庫)
	if err := HandleMessage(&message); err != nil {
		return err
	}
	// 取得該聊天室的成員列表
	userIDs, err := repository.GetUserIDsByChatID(chatID)
	if err != nil {
		return err
	}
	// 將訊息轉換為 JSON
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	// 廣播訊息給聊天室成員
	for _, uid := range userIDs {
		ws.SendToUser(uid, data)
	}

	return nil
}

// GetMessages 取得聊天室的分頁訊息
func GetMessages(chatID, lastID uint, limit int) ([]domain.Message, error) {
	return repository.GetMessagesByChatID(chatID, lastID, limit)
}
