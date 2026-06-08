package handler

import (
	"go-chat/internal/middleware"
	"go-chat/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMessagesHandler(c *gin.Context) {
	// 從 URL 參數取得 chatID
	chatIDStr := c.Param("id")
	// 從 get 參數取得 last_id（pagination）
	lastIDStr := c.Query("last_id")
	// chatID 轉換為 uint
	chatIDUint64, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	chatID := uint(chatIDUint64)
	// 取得 userID (由 middleware 從 JWT 解析)
	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// lastID 轉換為 uint
	var lastID uint64
	if lastIDStr != "" {
		lastID, err = strconv.ParseUint(lastIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid last_id"})
			return
		}
	}
	// 取得訊息列表
	msgs, err := service.GetMessages(chatID, uint(lastID), 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}
	// 自動標記已讀
	if len(msgs) > 0 && lastID == 0 {
		lastMsgID := msgs[0].ID
		_ = service.MarkMessagesRead(userID, chatID, lastMsgID)
	}
	// 回傳訊息列表
	c.JSON(http.StatusOK, msgs)
}
