package handler

import (
	"go-chat/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateGroupChatRequest struct {
	Name    string `json:"name"`
	UserIDs []uint `json:"user_ids"`
}

func CreateGroupChatHandler(c *gin.Context) {
	var req CreateGroupChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 假設 user_id 從 query 拿（之後可改 JWT）
	userIDStr := c.Query("user_id")
	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	chatID, err := service.CreateGroupChat(uint(userIDUint64), req.Name, req.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chat_id": chatID,
	})
}
