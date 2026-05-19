package handler

import (
	"errors"
	"go-chat/internal/middleware"
	"go-chat/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateGroupChatRequest struct {
	Name    string `json:"name"`
	UserIDs []uint `json:"user_ids"`
}

// 建立群組聊天室
func CreateGroupChatHandler(c *gin.Context) {
	var req CreateGroupChatRequest

	// ✅ 解析 request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// ✅ 驗證 name
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	currentUserID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// ✅ 去重 + 自動加入自己
	userMap := make(map[uint]bool)
	userMap[currentUserID] = true

	for _, uid := range req.UserIDs {
		userMap[uid] = true
	}

	var finalUserIDs []uint
	for uid := range userMap {
		finalUserIDs = append(finalUserIDs, uid)
	}

	// ✅ 呼叫 service
	chatID, err := service.CreateGroupChat(currentUserID, req.Name, finalUserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}

	// ✅ 回傳結果
	c.JSON(http.StatusCreated, gin.H{
		"chat_id": chatID,
	})
}
func GetChatsHandler(c *gin.Context) {
	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	chats, err := service.GetUserChats(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get chats"})
		return
	}

	c.JSON(200, chats)
}

type CreatePrivateChatRequest struct {
	TargetUserID uint `json:"target_user_id"`
}

func MarkReadHandler(c *gin.Context) {
	chatIDStr := c.Param("id")
	chatIDUint64, err := strconv.ParseUint(chatIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		LastReadMessageID uint `json:"last_read_message_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.LastReadMessageID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "last_read_message_id is required"})
		return
	}

	if err := service.MarkMessagesRead(userID, uint(chatIDUint64), req.LastReadMessageID); err != nil {
		if errors.Is(err, service.ErrNotChatMember) || errors.Is(err, service.ErrInvalidMessageID) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark messages as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func CreatePrivateChatHandler(c *gin.Context) {
	var req CreatePrivateChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	chatID, err := service.CreatePrivateChat(userID, req.TargetUserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"chat_id": chatID,
	})
}
