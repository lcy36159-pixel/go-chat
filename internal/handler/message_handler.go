package handler

import (
	"go-chat/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMessagesHandler(c *gin.Context) {
	user1 := c.Query("user1")
	user2 := c.Query("user2")
	lastIDStr := c.Query("last_id")

	var lastID uint64 = 0
	if lastIDStr != "" {
		lastID, _ = strconv.ParseUint(lastIDStr, 10, 64)
	}

	msgs, _ := repository.GetMessages(user1, user2, uint(lastID))

	c.JSON(http.StatusOK, msgs)
}
