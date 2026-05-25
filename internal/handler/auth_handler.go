package handler

import (
	"net/http"
	"strconv"

	"go-chat/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	userIDStr := c.Query("user_id")

	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	token, err := jwt.GenerateToken(uint(userIDUint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
