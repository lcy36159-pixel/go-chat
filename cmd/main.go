package main

import (
	"log"
	"os"

	"go-chat/internal/handler"
	"go-chat/internal/middleware"
	"go-chat/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to load .env: %v", err)
	}

	db.Init()

	// db.DB.AutoMigrate(
	// 	&domain.User{},
	// 	&domain.Chat{},
	// 	&domain.ChatMember{},
	// 	&domain.Message{},
	// 	&domain.MessageRead{},
	// )

	r := gin.Default()
	auth := r.Group("/")
	// login 不需要 auth
	r.GET("/login", handler.LoginHandler)
	// WebSocket 連線
	r.GET("/ws", handler.WebSocketHandler)
	// 需要登入的
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/chats/private", handler.CreatePrivateChatHandler)
	auth.POST("/chats/group", handler.CreateGroupChatHandler)
	auth.POST("/chats/:id/read", handler.MarkReadHandler)
	auth.GET("/messages", handler.GetMessagesHandler)
	auth.GET("/chats", handler.GetChatsHandler)
	if err := r.Run(":" + serverPort()); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func serverPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		return "8080"
	}
	return port
}
