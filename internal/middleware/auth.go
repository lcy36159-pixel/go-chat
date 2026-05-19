package middleware

import (
	"errors"
	"net/http"
	"strings"

	"go-chat/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const UserIDContextKey = "user_id"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		userID, err := jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(UserIDContextKey, userID)
		c.Next()
	}
}

func UserIDFromContext(c *gin.Context) (uint, error) {
	val, ok := c.Get(UserIDContextKey)
	if !ok {
		return 0, errors.New("missing authenticated user")
	}

	userID, ok := val.(uint)
	if !ok || userID == 0 {
		return 0, errors.New("invalid authenticated user")
	}

	return userID, nil
}
