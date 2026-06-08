package handler

import (
	"go-chat/internal/middleware"
	"go-chat/internal/repository"
	"go-chat/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMessagesHandler(c *gin.Context) {
	chatIDStr := c.Param("id")
	lastIDStr := c.Query("last_id")

	// 解析 id
	chatIDUint64, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	chatID := uint(chatIDUint64)

	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 解析 last_id（pagination）
	var lastID uint64
	if lastIDStr != "" {
		lastID, err = strconv.ParseUint(lastIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid last_id"})
			return
		}
	}

	// 取得訊息
	msgs, err := repository.GetMessagesByChatID(chatID, uint(lastID), 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	// 自動標記已讀
	if len(msgs) > 0 && lastID == 0 {
		lastMsgID := msgs[0].ID
		_ = service.MarkMessagesRead(userID, chatID, lastMsgID)
	}

	c.JSON(http.StatusOK, msgs)
}
