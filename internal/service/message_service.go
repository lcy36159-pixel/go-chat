package service

import (
	"encoding/json"
	"go-chat/internal/domain"
	"go-chat/internal/repository"
	"go-chat/internal/ws"
)

// SendMessage validates input, checks chat membership, saves the message,
// and broadcasts it to all members of the chat room.
func SendMessage(userID, chatID uint, content string) error {
	if userID == 0 {
		return ErrUserIDRequired
	}
	if chatID == 0 {
		return ErrChatIDRequired
	}
	if content == "" {
		return ErrContentEmpty
	}

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

	if err := HandleMessage(&message); err != nil {
		return err
	}

	userIDs, err := repository.GetUserIDsByChatID(chatID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for _, uid := range userIDs {
		ws.SendToUser(uid, data)
	}

	return nil
}

// GetMessages retrieves paginated messages for a chat room.
func GetMessages(chatID, lastID uint, limit int) ([]domain.Message, error) {
	return repository.GetMessagesByChatID(chatID, lastID, limit)
}
