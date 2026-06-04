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
	r.POST("/register", handler.RegisterHandler)
	r.POST("/login", handler.LoginHandler)
	// WebSocket 連線
	r.GET("/ws", handler.WebSocketHandler)
	// 需要登入的
	auth.Use(middleware.AuthMiddleware())
	// 建立私人聊天室
	auth.POST("/chats/private", handler.CreatePrivateChatHandler)
	// 建立群組聊天室
	auth.POST("/chats/group", handler.CreateGroupChatHandler)
	// 加入群組成員
	auth.POST("/chats/:id/members", handler.AddGroupMemberHandler)
	// 標記已讀
	auth.POST("/chats/:id/read", handler.MarkReadHandler)
	// 取得訊息
	auth.GET("/messages", handler.GetMessagesHandler)
	// 取得聊天室清單
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
