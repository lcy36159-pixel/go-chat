package handler

import (
	"go-chat/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMessagesHandler(c *gin.Context) {
	chatIDStr := c.Query("chat_id")
	lastIDStr := c.Query("last_id")

	// 解析 chat_id
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	// 解析 last_id（pagination）
	var lastID uint64 = 0
	if lastIDStr != "" {
		lastID, err = strconv.ParseUint(lastIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid last_id"})
			return
		}
	}

	msgs, err := repository.GetMessagesByChatID(uint(chatID), uint(lastID), 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, msgs)
}
