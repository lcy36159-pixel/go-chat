package middleware

import (
	"errors"
	"net/http"
	"strings"

	"go-chat/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const UserIDContextKey = "user_id"

// AuthMiddleware 驗證 JWT 並將 userID 存入 context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Authorization header 取得 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		// 驗證 token 開頭是否為 "Bearer "
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		// 解析 token 並取得 userID
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		userID, err := jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// 將解析出來的 userID 存入 context 中的 "user_id" 欄位
		c.Set(UserIDContextKey, userID)
		c.Next()
	}
}

// UserIDFromContext 從 context 中取得已驗證的使用者 ID
func UserIDFromContext(c *gin.Context) (uint, error) {
	// 從 context 中取得 userID，並確保它是 uint 類型且不為 0
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
