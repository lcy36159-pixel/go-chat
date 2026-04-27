package service

import (
	"errors"
	"go-chat/internal/domain"
	"go-chat/internal/repository"
	"go-chat/pkg/db"
)

//
// ========================
// 📩 Message 處理
// ========================
//

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

	// 預設 type
	if msg.Type == "" {
		msg.Type = "text"
	}

	return repository.SaveMessage(msg)
}

//
// ========================
// 💬 Chat 建立
// ========================
//

// 建立群組聊天室（完整安全版）
func CreateGroupChat(creatorID uint, name string, userIDs []uint) (uint, error) {
	if name == "" {
		return 0, errors.New("name is required")
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 建 chat
	chat := domain.Chat{
		Name:      name,
		Type:      "group",
		CreatedBy: creatorID,
	}

	if err := tx.Create(&chat).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// member 去重
	memberMap := make(map[uint]bool)
	memberMap[creatorID] = true

	for _, id := range userIDs {
		memberMap[id] = true
	}

	// 寫入 chat_members
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

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return chat.ID, nil
}

// MarkMessagesRead records the last read message for a user in a chat,
// updating their unread count cursor.
func MarkMessagesRead(userID, chatID, lastReadMessageID uint) error {
	if userID == 0 {
		return errors.New("user_id is required")
	}
	if chatID == 0 {
		return errors.New("chat_id is required")
	}
	if lastReadMessageID == 0 {
		return errors.New("last_read_message_id is required")
	}
	return repository.UpsertMessageRead(userID, chatID, lastReadMessageID)
}

//
// ========================
//

// 取得使用者聊天室列表
func GetUserChats(userID uint) ([]domain.ChatDTO, error) {
	return repository.GetUserChats(userID)
}

func CreatePrivateChat(user1 uint, user2 uint) (uint, error) {
	if user1 == user2 {
		return 0, errors.New("cannot chat with yourself")
	}

	// 🔍 先檢查是否已存在
	existingChatID, err := repository.FindPrivateChat(user1, user2)
	if err == nil && existingChatID != 0 {
		return existingChatID, nil
	}

	// 建立 chat
	key := repository.GeneratePrivateKey(user1, user2)
	chat := domain.Chat{
		Type:      "private",
		SearchKey: key,
	}

	if err := repository.CreateChat(&chat); err != nil {
		return 0, err
	}

	// 加入兩人
	repository.AddChatMember(&domain.ChatMember{
		ChatID: chat.ID,
		UserID: user1,
	})
	repository.AddChatMember(&domain.ChatMember{
		ChatID: chat.ID,
		UserID: user2,
	})

	return chat.ID, nil
}
