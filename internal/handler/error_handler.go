package handler

import (
	"errors"
	"go-chat/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError maps known service sentinel errors to appropriate HTTP responses.
// Unknown errors return HTTP 500 with a generic message.
func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotChatMember):
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not a member of this group"})
	case errors.Is(err, service.ErrNotGroupChat):
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat is not a group"})
	case errors.Is(err, service.ErrInvalidMessageID):
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid message"})
	case errors.Is(err, service.ErrAlreadyMember):
		c.JSON(http.StatusConflict, gin.H{"error": "user is already a member"})
	case errors.Is(err, service.ErrInvalidRegisterInput), errors.Is(err, service.ErrWeakPassword):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrAccountTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrInvalidLoginInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
